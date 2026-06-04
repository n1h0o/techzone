-- +goose Up

CREATE TABLE users(
    id BIGSERIAL PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'client',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose Down

DROP TABLE users
