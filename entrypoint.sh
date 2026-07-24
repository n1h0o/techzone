#!/bin/sh

set -e

if [ -z "${DB_URL}" ]; then
  DB_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:${DB_PORT}/${POSTGRES_DB}?sslmode=disable"
fi

echo "Running migrations..."
goose -dir migrations postgres "$DB_URL" up

echo "Starting server..."
exec ./server
