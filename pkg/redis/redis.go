package redis

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

// создает redis клиент из url или адреса и сразу проверяет соединение
func New() (*redis.Client, error) {
	var client *redis.Client

	if url := os.Getenv("REDIS_URL"); url != "" {
		// url удобен для облачных провайдеров которые выдают одну готовую строку
		log.Println("REDIS_URL =", os.Getenv("REDIS_URL"))
		opt, err := redis.ParseURL(url)
		if err != nil {
			return nil, err
		}
		client = redis.NewClient(opt)
	} else {
		// адрес по умолчанию удобен для локального docker compose
		addr := os.Getenv("REDIS_ADDR")
		if addr == "" {
			addr = "localhost:6379"
		}

		client = redis.NewClient(&redis.Options{
			Addr: addr,
		})
	}

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
