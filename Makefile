.PHONY: run build tidy fmt vet

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

migrate-up:
	goose -dir migrations postgres "$(DB_URL)" up
