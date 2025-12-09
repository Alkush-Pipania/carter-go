-- name: CreateFallbackImage :one
INSERT INTO fallback_images (imgurl)
VALUES ($1)
RETURNING *;

-- name: GetFallbackImageByID :one
SELECT * FROM fallback_images WHERE id = $1;

-- name: GetRandomFallbackImage :one
SELECT * FROM fallback_images ORDER BY RANDOM() LIMIT 1;

-- name: ListFallbackImages :many
SELECT * FROM fallback_images ORDER BY id;

-- name: DeleteFallbackImage :exec
DELETE FROM fallback_images WHERE id = $1;
