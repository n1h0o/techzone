package integration

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var db *pgxpool.Pool

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../.env"); err != nil {
		panic(err)
	}

	var err error

	db, err = pgxpool.New(
		context.Background(),
		os.Getenv("DB_URL"),
	)
	if err != nil {
		panic(err)
	}

	defer func() {
		if db != nil {
			db.Close()
		}
	}()

	code := m.Run()

	os.Exit(code)
}
