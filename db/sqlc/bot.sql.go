// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: bot.sql

package db

import (
	"context"
	"database/sql"
)

const createBot = `-- name: CreateBot :one
INSERT INTO bots (
    title,
    company_id
) VALUES (
    $1, $2
) RETURNING id, title, company_id, created_at, updated_at
`

type CreateBotParams struct {
	Title     string `json:"title"`
	CompanyID int64  `json:"company_id"`
}

func (q *Queries) CreateBot(ctx context.Context, arg CreateBotParams) (Bot, error) {
	row := q.db.QueryRowContext(ctx, createBot, arg.Title, arg.CompanyID)
	var i Bot
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.CompanyID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteBot = `-- name: DeleteBot :exec
DELETE FROM bots
WHERE id = $1
`

func (q *Queries) DeleteBot(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteBot, id)
	return err
}

const getBot = `-- name: GetBot :one
SELECT id, title, company_id, created_at, updated_at FROM bots
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetBot(ctx context.Context, id int64) (Bot, error) {
	row := q.db.QueryRowContext(ctx, getBot, id)
	var i Bot
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.CompanyID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listAllBots = `-- name: ListAllBots :many
SELECT id, title, company_id, created_at, updated_at FROM bots
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListAllBotsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAllBots(ctx context.Context, arg ListAllBotsParams) ([]Bot, error) {
	rows, err := q.db.QueryContext(ctx, listAllBots, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Bot{}
	for rows.Next() {
		var i Bot
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.CompanyID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listCompanyBots = `-- name: ListCompanyBots :many
SELECT id, title, company_id, created_at, updated_at FROM bots
WHERE company_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListCompanyBotsParams struct {
	CompanyID int64 `json:"company_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListCompanyBots(ctx context.Context, arg ListCompanyBotsParams) ([]Bot, error) {
	rows, err := q.db.QueryContext(ctx, listCompanyBots, arg.CompanyID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Bot{}
	for rows.Next() {
		var i Bot
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.CompanyID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateBot = `-- name: UpdateBot :one
UPDATE bots
SET
 title = coalesce($1, title),
 company_id = coalesce($2, company_id),
 updated_at = now()
WHERE id = $3 
RETURNING id, title, company_id, created_at, updated_at
`

type UpdateBotParams struct {
	Title     sql.NullString `json:"title"`
	CompanyID sql.NullInt64  `json:"company_id"`
	ID        int64          `json:"id"`
}

func (q *Queries) UpdateBot(ctx context.Context, arg UpdateBotParams) (Bot, error) {
	row := q.db.QueryRowContext(ctx, updateBot, arg.Title, arg.CompanyID, arg.ID)
	var i Bot
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.CompanyID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
