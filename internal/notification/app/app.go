package app

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"techzone/internal/config"
	grpcserver "techzone/internal/notification/grpc"
	"techzone/internal/notification/kafka"
	"techzone/internal/notification/repository"
	"techzone/internal/notification/service"
	"techzone/internal/notification/worker"
	pkg "techzone/pkg/kafka"
	"techzone/pkg/postgres"

	"google.golang.org/grpc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/twmb/franz-go/pkg/kgo"
	"time"
)

type App struct {
	// хранит долгоживущие зависимости notification сервиса
	grpcServer *grpc.Server

	config *config.Config

	notificationService *service.NotificationService

	db *pgxpool.Pool

	consumer       *kafka.Consumer
	consumerClient *kgo.Client

	cancel context.CancelFunc
}

// собирает зависимости notification сервиса отдельно от основного api
type Dependencies struct {
	Config         *config.Config
	DB             *pgxpool.Pool
	ConsumerClient *kgo.Client
	Cancel         context.CancelFunc
}

// поднимает зависимости notification сервиса
func BuildDependencies() (*Dependencies, error) {
	cfg := config.Load()

	db, err := postgres.New(cfg.DatabaseURL())
	if err != nil {
		return nil, err
	}

	consumerClient, err := pkg.NewConsumerClient()
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Dependencies{
		Config:         cfg,
		DB:             db,
		ConsumerClient: consumerClient,
	}, nil
}

// связывает зависимости с grpc сервером и kafka consumer
func New(
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
	notificationRepo := repository.NewNotificationRepository(deps.DB)
	notificationService := service.NewNotificationService(notificationRepo)

	workerPool := worker.NewNotificationWorkerPool(
		5,
		notificationService,
	)

	consumer := kafka.NewConsumer(
		deps.ConsumerClient,
		workerPool,
	)

	grpcServer := grpc.NewServer()

	grpcNotificationServer := grpcserver.NewServer(
		notificationService,
	)

	grpcserver.Register(
		grpcServer,
		grpcNotificationServer,
	)

	return &App{
		grpcServer:          grpcServer,
		config:              deps.Config,
		notificationService: notificationService,
		consumer:            consumer,
		db:                  deps.DB,
		consumerClient:      deps.ConsumerClient,
		cancel:              deps.Cancel,
	}, nil
}

// освобождает долгоживущие ресурсы notification сервиса
func (a *App) Close() {

	if a.cancel != nil {
		a.cancel()
	}
	if a.consumerClient != nil {
		a.consumerClient.Close()
	}
	if a.db != nil {
		a.db.Close()
	}
}

// запускает grpc сервер и sidecar метрики
func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}
	log.Println("Kafka consumer started")
	go a.consumer.Start(ctx)

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		// отдельный http сервер для метрик не мешает grpc трафику
		metricsServer := &http.Server{
			Addr:              ":9091",
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		}

		log.Println("metrics server started on :9091")

		if err := metricsServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Println(err)
		}
	}()

	log.Println("gRPC server started on :50051")

	return a.grpcServer.Serve(lis)
}
