-- name: CreateCompany :one
INSERT INTO companies (
    phone,
    name,
    email
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetCompanyByEmail :one
SELECT * FROM companies
WHERE email = $1 LIMIT 1;

-- name: GetCompanyByID :one
SELECT * FROM companies
WHERE id = $1 LIMIT 1;


-- name: UpdateCompany :one
UPDATE companies
SET
 email = coalesce(sqlc.narg('email'), email),
 phone = coalesce(sqlc.narg('phone'), phone),
 name = coalesce(sqlc.narg('name'), name),
 updated_at = now()
WHERE id = sqlc.arg('id') 
RETURNING *;


-- name: ListCompanies :many
SELECT * FROM companies
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteCompany :exec
DELETE FROM companies
WHERE id = $1;


