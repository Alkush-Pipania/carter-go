-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;