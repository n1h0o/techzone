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

	cartRepo := repository.NewCartRepository(db)
	cartService := service.NewCartService(cartRepo)
	cartHandler := handler.NewCartHandler(cartService)

	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, cartRepo, productRepo, db)
	orderHandler := handler.NewOrderHandler(orderService)

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

	handlerWithCors := middleware.CORSMiddleware(mux)

	log.Println("server started on :8080")

	if err := http.ListenAndServe(":8080", handlerWithCors); err != nil {
		log.Fatal(err)
	}
}
