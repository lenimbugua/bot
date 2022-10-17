-- name: CreateUser :one
INSERT INTO users (
  name,
  password_hash,
  phone,
  company_id
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE phone = $1 LIMIT 1;
