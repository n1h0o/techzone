package config

import "os"

type Config struct {
	DBURL     string
	JWTSecret string
}

func Load() *Config {
	return &Config{
		DBURL:     os.Getenv("DB_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
