.PHONY: run build tidy fmt vet test lint migrate-up migrate-down start stop restart

-include .env
export

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

unit:
	go test -race ./internal/service/... ./internal/repository/...

integration:
	go test ./internal/tests/integration -v

lint:
	golangci-lint run ./...

migrate-up:
	goose -dir migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DB_URL)" down

start:
	docker compose up -d --build
	cd frontend && npm run dev

stop:
	docker compose down

restart:
	docker compose down
	docker compose up -d --build
	cd frontend && npm run dev