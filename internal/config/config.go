package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost string
	DBPort string

	DBUser               string
	DBPassword           string
	DBName               string
	JWTSecret            string
	NotificationGRPCAddr string
}

func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found")
	}

	return &Config{

		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),

		DBUser:     os.Getenv("POSTGRES_USER"),
		DBPassword: os.Getenv("POSTGRES_PASSWORD"),
		DBName:     os.Getenv("POSTGRES_DB"),

		JWTSecret: os.Getenv("JWT_SECRET"),

		NotificationGRPCAddr: os.Getenv("NOTIFICATION_GRPC_ADDR"),
	}
}
