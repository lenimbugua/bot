// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: session.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createSession = `-- name: CreateSession :one
INSERT INTO sessions (
  id,
  user_id,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  channel_id,
  question_id,
  response_id,
  expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING id, user_id, refresh_token, user_agent, client_ip, is_blocked, channel_id, question_id, response_id, expires_at, created_at, updated_at
`

type CreateSessionParams struct {
	ID           uuid.UUID `json:"id"`
	UserID       int64     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ChannelID    int64     `json:"channel_id"`
	QuestionID   int64     `json:"question_id"`
	ResponseID   int64     `json:"response_id"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, createSession,
		arg.ID,
		arg.UserID,
		arg.RefreshToken,
		arg.UserAgent,
		arg.ClientIp,
		arg.IsBlocked,
		arg.ChannelID,
		arg.QuestionID,
		arg.ResponseID,
		arg.ExpiresAt,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ChannelID,
		&i.QuestionID,
		&i.ResponseID,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSession = `-- name: GetSession :one
SELECT id, user_id, refresh_token, user_agent, client_ip, is_blocked, channel_id, question_id, response_id, expires_at, created_at, updated_at FROM sessions
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetSession(ctx context.Context, id uuid.UUID) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSession, id)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ChannelID,
		&i.QuestionID,
		&i.ResponseID,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}