-- name: GetChatMessagesBySessionID :many
SELECT * FROM chat_messages WHERE session_id = $1;
