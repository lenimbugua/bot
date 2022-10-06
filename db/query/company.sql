-- name: CreateCompany :one
INSERT INTO companies (
    phone,
    name,
    email
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetCompany :one
SELECT * FROM companies
WHERE email = $1 LIMIT 1;
