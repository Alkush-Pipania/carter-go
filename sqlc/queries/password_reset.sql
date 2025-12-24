-- name: CreatePasswordResetToken :one
INSERT INTO password_reset_tokens (
  id, user_id, token_hash, expires_at
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetValidPasswordResetToken :one
SELECT * FROM password_reset_tokens
WHERE token_hash = $1
  AND expires_at > NOW()
  AND used_at IS NULL
LIMIT 1;

-- name: MarkPasswordResetTokenUsed :exec
UPDATE password_reset_tokens
SET used_at = NOW()
WHERE id = $1;
