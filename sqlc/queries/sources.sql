-- name: GetSourcesByUserID :many
SELECT * FROM sources WHERE user_id = $1;
