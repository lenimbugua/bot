-- name: CreateChannel :one
INSERT INTO channels (
    name
) VALUES (
    $1
) RETURNING *;

-- name: GetChannel :one
SELECT * FROM channels
WHERE name = $1 LIMIT 1;
