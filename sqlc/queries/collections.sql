-- name: GetCollectionsByUserID :many
SELECT * FROM collections WHERE user_id = $1;
