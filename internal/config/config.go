package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// хранит готовую строку подключения если она приходит целиком
	DBURL string

	// хранит составные части подключения когда строка подключения не передана
	DBHost string
	DBPort string

	// хранит секреты и адреса внешних зависимостей
	DBUser               string
	DBPassword           string
	DBName               string
	JWTSecret            string
	NotificationGRPCAddr string
}

// загружает настройки из окружения и локального файла для разработки
func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found")
	}

	return &Config{
		DBURL: os.Getenv("DB_URL"),

		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),

		DBUser:     os.Getenv("POSTGRES_USER"),
		DBPassword: os.Getenv("POSTGRES_PASSWORD"),
		DBName:     os.Getenv("POSTGRES_DB"),

		JWTSecret: os.Getenv("JWT_SECRET"),

		NotificationGRPCAddr: os.Getenv("NOTIFICATION_GRPC_ADDR"),
	}
}

// собирает строку подключения из частей если целое значение не задано
func (c *Config) DatabaseURL() string {
	if c.DBURL != "" {
		return c.DBURL
	}

	host := c.DBHost
	if host == "" {
		host = "localhost"
	}

	port := c.DBPort
	if port == "" {
		port = "5432"
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser,
		c.DBPassword,
		host,
		port,
		c.DBName,
	)
}
