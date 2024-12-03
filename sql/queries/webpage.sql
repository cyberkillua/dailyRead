-- name: CreateWebpage :one
INSERT INTO webpages (id, created_at, updated_at, name, url, type)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;


-- name: GetNextWebpageToFetch :many
SELECT * FROM webpages
ORDER BY last_updated_at ASC NULLS FIRST   
LIMIT $1;


-- name: MarkWebpageAsFetched :one
UPDATE webpages
SET last_updated_at = Now(), 
updated_at = Now()
WHERE id = $1
RETURNING *;