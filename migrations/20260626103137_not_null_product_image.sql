-- +goose Up
UPDATE products
SET image_url = ''
WHERE image_url IS NULL;

ALTER TABLE products
ALTER COLUMN image_url SET DEFAULT '';

ALTER TABLE products
ALTER COLUMN image_url SET NOT NULL;

-- +goose Down
ALTER TABLE products
ALTER COLUMN image_url DROP NOT NULL;

ALTER TABLE products
ALTER COLUMN image_url DROP DEFAULT;