// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import (
	"context"
)

type Querier interface {
	CreateCompany(ctx context.Context, arg CreateCompanyParams) (Company, error)
	GetCompany(ctx context.Context, email string) (Company, error)
}

var _ Querier = (*Queries)(nil)
