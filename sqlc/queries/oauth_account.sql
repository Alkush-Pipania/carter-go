-- name: GetOAuthAccount :one
SELECT * FROM oauth_accounts
WHERE provider = $1 AND provider_user_id = $2;

-- name: CreateOAuthAccount :one
INSERT INTO oauth_accounts (
  id, user_id, provider, provider_user_id, email
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
