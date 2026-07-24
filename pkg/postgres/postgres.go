package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// создает пул и несколько раз пингует базу пока контейнер не станет готов
func New(dbURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}

	for i := 0; i < 10; i++ {
		if err = pool.Ping(context.Background()); err == nil {
			return pool, nil
		}

		time.Sleep(500 * time.Millisecond)
	}

	pool.Close()

	return nil, fmt.Errorf("postgres connection failed: %v", err)
}
