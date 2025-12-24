-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (id, email, password_hash, verified)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2
WHERE id = $1;