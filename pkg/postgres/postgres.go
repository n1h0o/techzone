package postgres

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New() *pgxpool.Pool {

	dbURL := os.Getenv("DB_URL")

	log.Println("DB_URL =", dbURL)
	log.Println("DATABASE_URL =", os.Getenv("DATABASE_URL"))

	pool, err := pgxpool.New(
		context.Background(),
		dbURL,
	)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		err = pool.Ping(context.Background())

		if err == nil {
			log.Println("postgres connected")
			return pool
		}
		log.Printf("waiting postgres...(%d/10)", i+1)
		time.Sleep(500 * time.Millisecond)
	}
	log.Fatalf("postgres connection failed: %v", err)
	return nil
}
