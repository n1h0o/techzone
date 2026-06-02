package postgres

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func New() *pgxpool.Pool {
	err := godotenv.Load()

	dbURL := os.Getenv("DB_URL")

	log.Println("DB_URL =", dbURL)

	pool, err := pgxpool.New(
		context.Background(),
		dbURL,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}
	log.Println("postgres connected")
	return pool
}
