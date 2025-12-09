-- name: CreateVerification :one
INSERT INTO verifications (email, token, expires)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetVerificationByEmailAndToken :one
SELECT * FROM verifications
WHERE email = $1 AND token = $2 AND expires > NOW();

-- name: DeleteVerification :exec
DELETE FROM verifications WHERE id = $1;

-- name: DeleteVerificationByEmail :exec
DELETE FROM verifications WHERE email = $1;

-- name: DeleteExpiredVerifications :exec
DELETE FROM verifications WHERE expires <= NOW();

-- name: CreateForgotPassword :one
INSERT INTO forgot_passwords (email, token, expires)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetForgotPasswordByEmailAndToken :one
SELECT * FROM forgot_passwords
WHERE email = $1 AND token = $2 AND expires > NOW();

-- name: DeleteForgotPassword :exec
DELETE FROM forgot_passwords WHERE id = $1;

-- name: DeleteForgotPasswordByEmail :exec
DELETE FROM forgot_passwords WHERE email = $1;

-- name: DeleteExpiredForgotPasswords :exec
DELETE FROM forgot_passwords WHERE expires <= NOW();
