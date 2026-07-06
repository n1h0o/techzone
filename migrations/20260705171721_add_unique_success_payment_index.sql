-- +goose Up
CREATE UNIQUE INDEX unique_success_payment_per_order
ON payments(order_id)
WHERE status = 'success';

-- +goose Down
DROP INDEX unique_success_payment_per_order;
