-- name: CreateBot :one
INSERT INTO bots (
    title,
    company_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetBot :one
SELECT * FROM bots
WHERE id = $1 LIMIT 1;


-- name: ListAllBots :many
SELECT * FROM bots
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListCompanyBots :many
SELECT * FROM bots
WHERE company_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateBot :one
UPDATE bots
SET
 title = coalesce(sqlc.narg('title'), title),
 company_id = coalesce(sqlc.narg('company_id'), company_id),
 updated_at = now()
WHERE id = sqlc.arg('id') 
RETURNING *;


-- name: DeleteBot :exec
DELETE FROM bots
WHERE id = $1;
