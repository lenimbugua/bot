// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  name,
  password_hash,
  phone,
  company_id
) VALUES (
  $1, $2, $3, $4
) RETURNING id, phone, company_id, password_hash, password_changed_at, name, created_at, updated_at
`

type CreateUserParams struct {
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"`
	Phone        string `json:"phone"`
	CompanyID    int64  `json:"company_id"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Name,
		arg.PasswordHash,
		arg.Phone,
		arg.CompanyID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Phone,
		&i.CompanyID,
		&i.PasswordHash,
		&i.PasswordChangedAt,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, phone, company_id, password_hash, password_changed_at, name, created_at, updated_at FROM users
WHERE phone = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, phone string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, phone)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Phone,
		&i.CompanyID,
		&i.PasswordHash,
		&i.PasswordChangedAt,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
