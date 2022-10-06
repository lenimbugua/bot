package db

import (
	"context"
	"database/sql"
	"fmt"
)

//Provides all the fuctions to execute database queries as well as transactions
type Store interface {
	Querier
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

//execTx executes a function within a database transaction
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
