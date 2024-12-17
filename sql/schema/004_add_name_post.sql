-- +goose Up
ALTER TABLE posts ADD COLUMN postName TEXT;

-- +goose Down
ALTER TABLE posts DROP COLUMN postName;