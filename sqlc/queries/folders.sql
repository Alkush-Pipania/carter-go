-- name: CreateFolder :one
INSERT INTO folders (name, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetFolderByID :one
SELECT * FROM folders WHERE id = $1 AND is_deleted = FALSE;

-- name: GetFolderBySecretKey :one
SELECT * FROM folders WHERE secret_key = $1 AND is_deleted = FALSE;

-- name: GetFoldersByUserID :many
SELECT * FROM folders
WHERE user_id = $1 AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: UpdateFolder :one
UPDATE folders
SET name = $2, updated_at = NOW()
WHERE id = $1 AND is_deleted = FALSE
RETURNING *;

-- name: SoftDeleteFolder :exec
UPDATE folders
SET is_deleted = TRUE, deleted_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: RestoreFolder :exec
UPDATE folders
SET is_deleted = FALSE, deleted_at = NULL, updated_at = NOW()
WHERE id = $1;

-- name: HardDeleteFolder :exec
DELETE FROM folders WHERE id = $1;

-- name: GetDeletedFoldersByUserID :many
SELECT * FROM folders
WHERE user_id = $1 AND is_deleted = TRUE
ORDER BY deleted_at DESC;
