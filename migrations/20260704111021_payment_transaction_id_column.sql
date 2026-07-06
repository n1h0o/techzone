-- +goose Up
ALTER TABLE payments
ADD COLUMN transaction_id TEXT;

-- +goose Down
ALTER TABLE payments
DROP COLUMN transaction_id;
