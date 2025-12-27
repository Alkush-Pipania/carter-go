-- name: GetSourceContentBySourceID :many
SELECT * FROM source_contents WHERE source_id = $1;

-- name: CreateSourceContent :one
INSERT INTO source_contents (source_id, content_text, content_hash)
VALUES ($1, $2, $3)
RETURNING *;
