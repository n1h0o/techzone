package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"techzone/internal/config"
	"techzone/internal/event"
	notificationgrpc "techzone/internal/grpcclient/notification"
	"techzone/internal/handler"
	"techzone/internal/middleware"
	"techzone/internal/payment"
	"techzone/internal/repository"
	"techzone/internal/seed"
	"techzone/internal/service"
	pkg "techzone/pkg/kafka"
	"techzone/pkg/postgres"
	"techzone/pkg/redis"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	goredis "github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/twmb/franz-go/pkg/kgo"
)

type MockPublisher struct{}

func (MockPublisher) PublishOrderCreated(
	context.Context,
	event.OrderCreatedEvent,
) error {
	return nil
}

func (MockPublisher) PublishPaymentCompleted(
	context.Context,
	event.PaymentCompletedEvent,
) error {
	return nil
}

type App struct {
	// хранит внешние клиенты чтобы их можно было закрыть при остановке
	notificationClient *notificationgrpc.Client

	handler http.Handler

	db *pgxpool.Pool

	redisClient *goredis.Client

	producerClient *kgo.Client
	consumerClient *kgo.Client

	cancel context.CancelFunc
}

// управляет тем как поднимается приложение в обычном режиме и в тестах
type ServerOptions struct {
	TestMode  bool
	SeedAdmin bool
}

// собирает зависимости до создания http маршрутов
type Dependencies struct {
	NotificationClient *notificationgrpc.Client
	Config             *config.Config
	DB                 *pgxpool.Pool
	RedisClient        *goredis.Client
	EventPublisher     service.EventPublisher
	PaymentGateway     payment.Gateway
	ProducerClient     *kgo.Client
	ConsumerClient     *kgo.Client
	Cancel             context.CancelFunc
	SeedAdmin          bool
}

// создает сервер с настройками по умолчанию
func NewServer(testMode bool) (*App, error) {
	return NewServerWithOptions(ServerOptions{
		TestMode:  testMode,
		SeedAdmin: true,
	})
}

// создает сервер с явными опциями для тестов и запуска
func NewServerWithOptions(
	options ServerOptions,
) (*App, error) {
	deps, err := BuildDependencies(options)
	if err != nil {
		return nil, err
	}

	app, err := NewWithDependencies(deps)
	if err != nil {
		closeDependencies(deps)
		return nil, err
	}

	return app, nil
}

// поднимает базовые зависимости приложения
func BuildDependencies(
	options ServerOptions,
) (*Dependencies, error) {
	cfg := config.Load()

	db, err := postgres.New(cfg.DatabaseURL())
	if err != nil {
		return nil, err
	}

	redisClient, err := redis.New()
	if err != nil {
		db.Close()
		return nil, err
	}

	notificationClient, err := notificationgrpc.New(cfg.NotificationGRPCAddr)
	if err != nil {
		if err := redisClient.Close(); err != nil {
			log.Printf("failed to close redis client: %v", err)
		}
		db.Close()
		return nil, err
	}

	producerClient, err := pkg.NewProducerClient()
	if err != nil {
		if err := notificationClient.Close(); err != nil {
			log.Printf("failed to close notification client: %v", err)
		}

		if err := redisClient.Close(); err != nil {
			log.Printf("failed to close redis client: %v", err)
		}

		db.Close()
		return nil, err
	}

	publisher := pkg.NewProducer(producerClient)

	deps := &Dependencies{
		NotificationClient: notificationClient,
		Config:             cfg,
		DB:                 db,
		RedisClient:        redisClient,
		EventPublisher:     publisher,
		ProducerClient:     producerClient,
		PaymentGateway:     payment.NewMockGateway(),
		SeedAdmin:          options.SeedAdmin,
	}

	if options.TestMode || os.Getenv("KAFKA_BROKERS") == "" {
		// в тестах и локальных режимах без kafka сервер поднимается без consumer части
		return deps, nil
	}

	return deps, nil
}

// связывает зависимости со слоями приложения и маршрутизатором
func NewWithDependencies(
	deps *Dependencies,
) (*App, error) {
	if deps == nil {
		return nil, errors.New("dependencies are nil")
	}
	if deps.Config == nil {
		return nil, errors.New("config is nil")
	}
	if deps.DB == nil {
		return nil, errors.New("db is nil")
	}
	if deps.RedisClient == nil {
		return nil, errors.New("redis client is nil")
	}
	if deps.EventPublisher == nil {
		return nil, errors.New("event publisher is nil")
	}
	if deps.PaymentGateway == nil {
		return nil, errors.New("payment gateway is nil")
	}
	if deps.NotificationClient == nil {
		return nil, errors.New("notification client is nil")
	}

	userRepo := repository.NewUserRepository(deps.DB)
	if deps.SeedAdmin {
		if err := seed.CreateAdmin(userRepo); err != nil {
			return nil, err
		}
	}

	notificationHandler := handler.NewNotificationHandler(deps.NotificationClient)

	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService, deps.Config)

	productRepo := repository.NewProductRepository(deps.DB)
	productService := service.NewProductService(productRepo, deps.RedisClient)
	productHandler := handler.NewProductHandler(productService)

	cartRepo := repository.NewCartRepository(deps.DB)
	cartService := service.NewCartService(cartRepo)
	cartHandler := handler.NewCartHandler(cartService)

	orderRepo := repository.NewOrderRepository(deps.DB)

	orderService := service.NewOrderService(
		orderRepo,
		cartRepo,
		productRepo,
		deps.EventPublisher,
		deps.DB,
		deps.RedisClient,
	)
	orderHandler := handler.NewOrderHandler(orderService)

	paymentService := service.NewPaymentService(
		deps.PaymentGateway,
		deps.EventPublisher,
		deps.DB,
	)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	mux := http.NewServeMux()

	// здесь концентрируется весь http контракт сервиса
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.Handle(
		"GET /me",
		middleware.AuthMiddleware(deps.Config)(
			http.HandlerFunc(
				handler.GetMe,
			),
		),
	)
	mux.HandleFunc("GET /products", productHandler.GetProducts)
	mux.Handle(
		"GET /admin/products",
		middleware.AuthMiddleware(deps.Config)(
			middleware.AdminMiddleware(
				http.HandlerFunc(
					productHandler.GetProductsForAdmin,
				),
			),
		),
	)
	mux.Handle(
		"POST /products",
		middleware.AuthMiddleware(deps.Config)(
			middleware.AdminMiddleware(
				http.HandlerFunc(
					productHandler.CreateProduct,
				),
			),
		),
	)
	mux.Handle(
		"PUT /products/{id}",
		middleware.AuthMiddleware(deps.Config)(
			middleware.AdminMiddleware(
				http.HandlerFunc(productHandler.UpdateProduct),
			),
		),
	)
	mux.Handle(
		"PATCH /products/{id}/status",
		middleware.AuthMiddleware(deps.Config)(
			middleware.AdminMiddleware(
				http.HandlerFunc(productHandler.SetProductStatus),
			),
		),
	)

	mux.HandleFunc("GET /products/{id}", productHandler.GetProduct)
	mux.Handle(
		"POST /cart/items",
		middleware.AuthMiddleware(deps.Config)(
			http.HandlerFunc(
				cartHandler.AddToCart,
			),
		),
	)
	mux.Handle(
		"GET /cart",
		middleware.AuthMiddleware(deps.Config)(
			http.HandlerFunc(cartHandler.GetCart),
		),
	)
	mux.Handle(
		"DELETE /cart/items/{item_id}",
		middleware.AuthMiddleware(deps.Config)(
			http.HandlerFunc(
				cartHandler.DeleteItem,
			),
		),
	)

	mux.Handle(
		"POST /orders",
		middleware.AuthMiddleware(deps.Config)(
			http.HandlerFunc(orderHandler.CreateOrder),
		),
	)
	mux.Handle(
		"GET /orders",
		middleware.AuthMiddleware(deps.Config)(
			http.HandlerFunc(orderHandler.GetOrders),
		),
	)
	mux.Handle(
		"GET /orders/{id}",
		middleware.AuthMiddleware(deps.Config)(
			http.HandlerFunc(orderHandler.GetOrderByID),
		),
	)

	mux.Handle(
		"PATCH /orders/{id}/status",
		middleware.AuthMiddleware(deps.Config)(
			http.HandlerFunc(orderHandler.UpdateStatus),
		),
	)

	mux.Handle(
		"GET /notifications",
		middleware.AuthMiddleware(deps.Config)(
			http.HandlerFunc(notificationHandler.GetNotifications),
		),
	)

	mux.Handle(
		"POST /payments",
		middleware.AuthMiddleware(deps.Config)(
			http.HandlerFunc(paymentHandler.Pay),
		),
	)

	mux.Handle(
		"GET /swagger/",
		httpSwagger.WrapHandler,
	)

	mux.Handle(
		"GET /metrics",
		promhttp.Handler(),
	)

	handler := middleware.MetricsMiddleware(mux)
	handler = middleware.CORSMiddleware(handler)

	return &App{

		notificationClient: deps.NotificationClient,

		handler: handler,

		db: deps.DB,

		redisClient: deps.RedisClient,

		producerClient: deps.ProducerClient,
		consumerClient: deps.ConsumerClient,

		cancel: deps.Cancel,
	}, nil

}

// закрывает зависимости если сборка сервера оборвалась на полпути
func closeDependencies(deps *Dependencies) {
	if deps == nil {
		return
	}

	if err := deps.NotificationClient.Close(); err != nil {
		log.Printf("failed to close notification client: %v", err)
	}

	if deps.Cancel != nil {
		deps.Cancel()
	}

	if deps.ConsumerClient != nil {
		deps.ConsumerClient.Close()
	}

	if deps.ProducerClient != nil {
		deps.ProducerClient.Close()
	}

	if deps.RedisClient != nil {
		if err := deps.RedisClient.Close(); err != nil {
			log.Printf("failed to close redis client: %v", err)
		}
	}

	if deps.DB != nil {
		deps.DB.Close()
	}
}

// возвращает готовый http обработчик для запуска сервера
func (a *App) Handler() http.Handler {
	return a.handler
}

// освобождает внешние ресурсы приложения
func (a *App) Close() {

	if a.cancel != nil {
		a.cancel()
	}

	if err := a.notificationClient.Close(); err != nil {
		log.Printf("failed to close notification client: %v", err)
	}

	if a.consumerClient != nil {
		a.consumerClient.Close()
	}

	if a.producerClient != nil {
		a.producerClient.Close()
	}

	if a.redisClient != nil {
		if err := a.redisClient.Close(); err != nil {
			log.Printf("failed to close redis client: %v", err)
		}
	}

	if a.db != nil {
		a.db.Close()
	}
}
