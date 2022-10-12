-- name: CreateUser :one
INSERT INTO users (
  name,
  password_hash,
  phone
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE phone = $1 LIMIT 1;
