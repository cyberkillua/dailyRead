-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, description, url, published_at, postName)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;


-- name: GetPosts :many
SELECT id, created_at, updated_at, title, description, url, published_at, postname 
FROM posts 
ORDER BY created_at DESC 
LIMIT 30;
