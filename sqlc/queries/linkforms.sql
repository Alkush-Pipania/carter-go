-- name: CreateLinkform :one
INSERT INTO linkforms (title, description, links, body, imgurl, tobefind, user_id, folder_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetLinkformBySecretID :one
SELECT * FROM linkforms WHERE secret_id = $1 AND is_deleted = FALSE;

-- name: GetLinkformsByUserID :many
SELECT * FROM linkforms
WHERE user_id = $1 AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: GetLinkformsByFolderID :many
SELECT * FROM linkforms
WHERE folder_id = $1 AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: GetPublicLinkforms :many
SELECT * FROM linkforms
WHERE tobefind = TRUE AND is_deleted = FALSE
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateLinkform :one
UPDATE linkforms
SET title = COALESCE($2, title),
    description = COALESCE($3, description),
    links = COALESCE($4, links),
    body = COALESCE($5, body),
    imgurl = COALESCE($6, imgurl),
    tobefind = COALESCE($7, tobefind),
    folder_id = $8
WHERE secret_id = $1 AND is_deleted = FALSE
RETURNING *;

-- name: SoftDeleteLinkform :exec
UPDATE linkforms
SET is_deleted = TRUE, deleted_at = NOW()
WHERE secret_id = $1;

-- name: RestoreLinkform :exec
UPDATE linkforms
SET is_deleted = FALSE, deleted_at = NULL
WHERE secret_id = $1;

-- name: HardDeleteLinkform :exec
DELETE FROM linkforms WHERE secret_id = $1;

-- name: GetDeletedLinkformsByUserID :many
SELECT * FROM linkforms
WHERE user_id = $1 AND is_deleted = TRUE
ORDER BY deleted_at DESC;

-- name: SearchLinkforms :many
SELECT * FROM linkforms
WHERE user_id = $1 
  AND is_deleted = FALSE
  AND (title ILIKE '%' || $2 || '%' OR description ILIKE '%' || $2 || '%')
ORDER BY created_at DESC;
