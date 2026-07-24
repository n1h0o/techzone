package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"techzone/internal/app"
	"techzone/internal/metrics"
	"time"

	_ "techzone/docs"

	"github.com/joho/godotenv"
)

// @title TechZone API
// @version 1.0
// @description Backend интернет-магазина TechZone
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// поднимает http сервер и завершает его через graceful shutdown
func main() {

	_ = godotenv.Load()

	application, err := app.NewServer(false)
	if err != nil {
		log.Fatal(err)
	}
	defer application.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		// ограничивает медленные соединения и защищает процесс от зависших клиентов
		Addr:              ":" + port,
		Handler:           application.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	metrics.Init()

	go func() {
		log.Printf("server started on :%s", port)

		if err := srv.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	log.Println("shutting down server")

	// дает активным запросам короткое окно чтобы завершиться без обрыва
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("server stopped")

}
