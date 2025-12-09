-- name: CreateUser :one
INSERT INTO users (email, username, image, password)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserBySecretKey :one
SELECT * FROM users WHERE secretkey = $1;

-- name: UpdateUser :one
UPDATE users
SET username = COALESCE($2, username),
    image = COALESCE($3, image),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2, updated_at = NOW()
WHERE id = $1;

-- name: VerifyUser :exec
UPDATE users
SET verified = TRUE, updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
