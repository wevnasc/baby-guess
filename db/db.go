package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Connection struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type Store struct {
	*sql.DB
}

func (store *Store) ExecTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := store.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	err = fn(tx)

	if err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, err)
		}
	}

	return tx.Commit()
}
