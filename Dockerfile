FROM golang:1.26 AS builder

WORKDIR /techzone

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /techzone/server .
COPY --from=builder /techzone/migrations ./migrations

EXPOSE 8080

CMD ["./server"]