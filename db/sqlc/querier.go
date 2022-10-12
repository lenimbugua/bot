// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateChannel(ctx context.Context, name string) (Channel, error)
	CreateCompany(ctx context.Context, arg CreateCompanyParams) (Company, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetChannel(ctx context.Context, name string) (Channel, error)
	GetCompany(ctx context.Context, email string) (Company, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetUser(ctx context.Context, phone string) (User, error)
}

var _ Querier = (*Queries)(nil)
