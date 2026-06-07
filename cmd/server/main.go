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

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.GetHealth)
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
		"POST /products",
		middleware.AuthMiddleware(cfg)(
			middleware.AdminMiddleware(
				http.HandlerFunc(
					productHandler.CreateProduct,
				),
			),
		),
	)
	mux.HandleFunc("GET /products/{id}", productHandler.GetProduct)

	log.Println("server started on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
