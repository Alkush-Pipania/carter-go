-- name: GetSourcesByUserID :many
SELECT * FROM sources WHERE user_id = $1;

-- name: GetSourcesByCollectionID :many
SELECT * FROM sources WHERE collection_id = $1;

-- name: GetSourceByID :one
SELECT * FROM sources WHERE id = $1;

-- name: CreateSource :one
INSERT INTO sources (user_id, collection_id, type, title, original_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: CreateSourceWithS3Key :one
INSERT INTO sources (user_id, collection_id, type, title, s3_key)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteSource :exec
DELETE FROM sources WHERE id = $1;

-- name: UpdateSourceStatus :one
UPDATE sources SET status = $2 WHERE id = $1
RETURNING *;
