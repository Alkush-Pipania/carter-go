-- name: GetCollectionsByUserID :many
SELECT * FROM collections WHERE user_id = $1;

-- name: CreateCollection :one
INSERT INTO collections (user_id, name)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateCollection :one
UPDATE collections
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteCollection :exec
DELETE FROM collections
WHERE id = $1;

