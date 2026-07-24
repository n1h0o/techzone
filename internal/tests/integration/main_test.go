package integration

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var db *pgxpool.Pool
var integrationSkipReason string

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../../.env")

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		integrationSkipReason = "DB_URL is not set"
		code := m.Run()
		os.Exit(code)
	}

	var err error

	db, err = pgxpool.New(
		context.Background(),
		dbURL,
	)
	if err != nil {
		integrationSkipReason = err.Error()
		code := m.Run()
		os.Exit(code)
	}

	if err := db.Ping(context.Background()); err != nil {
		integrationSkipReason = err.Error()
		db.Close()
		db = nil
	}

	code := m.Run()

	if db != nil {
		db.Close()
	}

	os.Exit(code)
}

func requireIntegration(t *testing.T) {
	t.Helper()

	if integrationSkipReason != "" {
		t.Skipf("integration tests skipped: %s", integrationSkipReason)
	}
}
