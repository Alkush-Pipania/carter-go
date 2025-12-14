-- name: GetChatSessionsByUserID :many
SELECT * FROM chat_sessions WHERE user_id = $1;
