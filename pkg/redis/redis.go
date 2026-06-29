package redis

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

func New() (*redis.Client, error) {
	addr := os.Getenv("REDIS_ADDR")

	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
