package main

import (
	"log"
	"net/http"
	"techzone/internal/config"
	"techzone/internal/handler"
	"techzone/internal/middleware"
	"techzone/internal/repository"
	"techzone/internal/service"
	"techzone/pkg/postgres"
)

func main() {
	cfg := config.Load()
	db := postgres.New()
	defer db.Close()

	userRepo := repository.NewUserRepository(db)

	authService := service.NewAuthService(userRepo)

	authHandler := handler.NewAuthHandler(authService, cfg)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.GetHealth)
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.Handle("GET /me", middleware.AuthMiddleware(cfg)(http.HandlerFunc(handler.GetMe)))

	log.Println("server started on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
