-- name: CreateSession :one
INSERT INTO sessions (
  id, user_id, expires_at, user_agent, ip_address
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM sessions
WHERE id = $1 AND is_revoked = false;

-- name: UpdateSessionLastSeen :exec
UPDATE sessions
SET last_seen_at = NOW()
WHERE id = $1;

-- name: RevokeSession :exec
UPDATE sessions
SET is_revoked = true
WHERE id = $1;

-- name: RevokeAllUserSessions :exec
UPDATE sessions
SET is_revoked = true
WHERE user_id = $1;
