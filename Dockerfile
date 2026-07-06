FROM golang:1.26 AS builder

WORKDIR /techzone

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /techzone/server .
COPY --from=builder /techzone/migrations ./migrations

EXPOSE 8080

CMD ["./server"]