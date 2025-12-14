-- name: GetPasswordResetByToken :one
SELECT * FROM password_resets WHERE token = $1;
