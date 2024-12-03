-- name: createWebpage :one
INSERT INTO webpages (id, created_at, updated_at, name, url, type)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

