package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New() (*pgxpool.Pool, error) {

	dbURL := os.Getenv("DB_URL")

	pool, err := pgxpool.New(
		context.Background(),
		dbURL,
	)
	if err != nil {
		return nil, err
	}

	for i := 0; i < 10; i++ {
		err = pool.Ping(context.Background())

		if err == nil {
			return pool, nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	pool.Close()

	return nil, fmt.Errorf("postgres connection failed: %w", err)
}
