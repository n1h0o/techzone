package main

import (
	"log"
	"net/http"
	"techzone/internal/handler"
	"techzone/pkg/postgres"
)

func main() {
	db := postgres.New()
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.GetHealth)
	log.Println("server started on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
