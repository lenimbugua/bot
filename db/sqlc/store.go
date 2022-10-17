package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Provides all the fuctions to execute database queries as well as transactions
type Store interface {
	Querier
	// UserCompanyTx(ctx context.Context, arg UserCompanyTxParams) (UserCompanyTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewSQLStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (dbStore *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := dbStore.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// UserCompanyParams contains the input parameters of the company user transaction
// type UserCompanyTxParams struct {
// 	UserID       int64  `json:"user_id"`
// 	CompanyPhone string `json:"phone"`
// 	CompanyName  string `json:"name"`
// 	CompanyEmail string `json:"email"`
// }

// // UserCompanyTxResult is the result of the  user company transaction
// type UserCompanyTxResult struct {
// 	Company     Company     `json:"company"`
// 	UserCompany UserCompany `json:"user_company"`
// }

// // UserCompanyTx creates a company and adds user to company user
// func (store *SQLStore) UserCompanyTx(ctx context.Context, arg UserCompanyTxParams) (UserCompanyTxResult, error) {
// 	var result UserCompanyTxResult

// 	err := store.execTx(ctx, func(q *Queries) error {
// 		var err error

// 		result.Company, err = q.CreateCompany(ctx, CreateCompanyParams{
// 			Phone: arg.CompanyPhone,
// 			Name:  arg.CompanyName,
// 			Email: arg.CompanyEmail,
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		result.UserCompany, err = q.CreateUserCompany(ctx, CreateUserCompanyParams{
// 			CompanyID: result.Company.ID,
// 			UserID:    arg.UserID,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return err
// 	})

// 	return result, err
// }
