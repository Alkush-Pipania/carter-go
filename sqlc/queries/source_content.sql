-- name: GetSourceContentBySourceID :many
SELECT * FROM source_contents WHERE source_id = $1;
