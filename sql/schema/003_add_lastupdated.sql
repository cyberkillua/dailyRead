-- +goose Up
ALTER TABLE webpages ADD COLUMN last_updated_at TIMESTAMP;

-- +goose Down
ALTER TABLE webpages DROP COLUMN last_updated_at;