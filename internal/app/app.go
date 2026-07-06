package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"techzone/internal/config"
	"techzone/internal/event"
	"techzone/internal/handler"
	"techzone/internal/middleware"
	"techzone/internal/payment"
	"techzone/internal/repository"
	"techzone/internal/service"
	"techzone/pkg/kafka"
	"techzone/pkg/postgres"
	"techzone/pkg/redis"

	"github.com/jackc/pgx/v5/pgxpool"
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
	handler http.Handler

	db *pgxpool.Pool

	redisClient *goredis.Client

	producerClient *kgo.Client
	consumerClient *kgo.Client

	cancel context.CancelFunc
}

func NewServer(testMode bool) *App {

	cfg := config.Load()
	db := postgres.New()

	redisClient, err := redis.New()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)

	authService := service.NewAuthService(userRepo)

	authHandler := handler.NewAuthHandler(authService, cfg)

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo, redisClient)
	productHandler := handler.NewProductHandler(productService)

	cartRepo := repository.NewCartRepository(db)
	cartService := service.NewCartService(cartRepo)
	cartHandler := handler.NewCartHandler(cartService)

	orderRepo := repository.NewOrderRepository(db)
	notificationRepo := repository.NewNotificationRepository(
		db,
	)
	notificationService := service.NewNotificationService(
		notificationRepo,
	)
	notificationPool := service.NewNotificationWorkerPool(
		5,
		notificationService,
	)
	notificationHandler := handler.NewNotificationHandler(
		notificationService,
	)

	var producer service.EventPublisher

	var producerClient *kgo.Client
	var consumerClient *kgo.Client
	var cancel context.CancelFunc

	gateway := payment.NewMockGateway()

	if testMode || os.Getenv("KAFKA_BROKERS") == "" {

		producer = MockPublisher{}

	} else {

		producerClient, err = kafka.NewProducerClient()
		if err != nil {
			log.Fatal(err)
		}

		consumerClient, err = kafka.NewConsumerClient()
		if err != nil {
			log.Fatal(err)
		}

		producer = kafka.NewProducer(producerClient)

		consumer := kafka.NewConsumer(
			consumerClient,
			notificationPool,
		)

		ctx, c := context.WithCancel(context.Background())
		cancel = c

		go consumer.Start(ctx)
	}

	orderService := service.NewOrderService(orderRepo, cartRepo, productRepo, producer, db)
	orderHandler := handler.NewOrderHandler(orderService)

	paymentService := service.NewPaymentService(gateway, producer, db)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.Handle(
		"GET /me",
		middleware.AuthMiddleware(cfg)(
			http.HandlerFunc(
				handler.GetMe,
			),
		),
	)
	mux.HandleFunc("GET /products", productHandler.GetProducts)
	mux.Handle(
		"GET /admin/products",
		middleware.AuthMiddleware(cfg)(
			middleware.AdminMiddleware(
				http.HandlerFunc(
					productHandler.GetProductsForAdmin,
				),
			),
		),
	)
	mux.Handle(
		"POST /products",
		middleware.AuthMiddleware(cfg)(
			middleware.AdminMiddleware(
				http.HandlerFunc(
					productHandler.CreateProduct,
				),
			),
		),
	)
	mux.Handle(
		"PUT /products/{id}",
		middleware.AuthMiddleware(cfg)(
			middleware.AdminMiddleware(
				http.HandlerFunc(productHandler.UpdateProduct),
			),
		),
	)
	mux.Handle(
		"PATCH /products/{id}/status",
		middleware.AuthMiddleware(cfg)(
			middleware.AdminMiddleware(
				http.HandlerFunc(productHandler.SetProductStatus),
			),
		),
	)

	mux.HandleFunc("GET /products/{id}", productHandler.GetProduct)
	mux.Handle(
		"POST /cart/items",
		middleware.AuthMiddleware(cfg)(
			http.HandlerFunc(
				cartHandler.AddToCart,
			),
		),
	)
	mux.Handle(
		"GET /cart",
		middleware.AuthMiddleware(cfg)(
			http.HandlerFunc(cartHandler.GetCart),
		),
	)
	mux.Handle(
		"DELETE /cart/items/{item_id}",
		middleware.AuthMiddleware(cfg)(
			http.HandlerFunc(
				cartHandler.DeleteItem,
			),
		),
	)

	mux.Handle(
		"POST /orders",
		middleware.AuthMiddleware(cfg)(
			http.HandlerFunc(orderHandler.CreateOrder),
		),
	)
	mux.Handle(
		"GET /orders",
		middleware.AuthMiddleware(cfg)(
			http.HandlerFunc(orderHandler.GetOrders),
		),
	)
	mux.Handle(
		"GET /orders/{id}",
		middleware.AuthMiddleware(cfg)(
			http.HandlerFunc(orderHandler.GetOrderByID),
		),
	)

	mux.Handle(
		"PATCH /orders/{id}/status",
		middleware.AuthMiddleware(cfg)(
			http.HandlerFunc(orderHandler.UpdateStatus),
		),
	)

	mux.Handle(
		"GET /notifications",
		middleware.AuthMiddleware(cfg)(
			http.HandlerFunc(
				notificationHandler.GetNotifications,
			),
		),
	)

	mux.Handle(
		"POST /payments",
		middleware.AuthMiddleware(cfg)(
			http.HandlerFunc(paymentHandler.Pay),
		),
	)

	mux.Handle(
		"GET /swagger/",
		httpSwagger.WrapHandler,
	)

	handlerWithCors := middleware.CORSMiddleware(mux)

	return &App{
		handler: handlerWithCors,

		db: db,

		redisClient: redisClient,

		producerClient: producerClient,
		consumerClient: consumerClient,

		cancel: cancel,
	}

}

func (a *App) Handler() http.Handler {
	return a.handler
}

func (a *App) Close() {

	if a.cancel != nil {
		a.cancel()
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
