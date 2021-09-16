package db

import (
	"context"
	"database/sql"
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

	defer func() {
		if err != nil {
			err = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	if err != nil {
		return err
	}

	err = fn(tx)

	if err != nil {
		return err
	}

	return nil
}
