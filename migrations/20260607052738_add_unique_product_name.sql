-- +goose Up
ALTER TABLE products
ADD CONSTRAINT products_name_unique UNIQUE (name);

-- +goose Down
SELECT 'down SQL query';
ALTER TABLE products
DROP CONSTRAINT products_name_unique;
