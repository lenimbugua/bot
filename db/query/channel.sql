-- name: CreateChannel :one
INSERT INTO channels (
    name
) VALUES (
    $1
) RETURNING *;

-- name: GetChannel :one
SELECT * FROM channels
WHERE name = $1 LIMIT 1;

-- name: ListChannels :many
SELECT * FROM channels
ORDER BY id
LIMIT $1
OFFSET $2;


-- name: UpdateChannel :one
UPDATE channels
SET
 name = coalesce(sqlc.narg('name'), name),
 updated_at = now()
WHERE id = sqlc.arg('id') 
RETURNING *;


-- name: DeleteChannel :exec
DELETE FROM channels
WHERE id = $1;

