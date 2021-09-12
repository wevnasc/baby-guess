package tables

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/wevnasc/baby-guess/db"
)

type Database struct {
	DB *db.Store
}

func newDatabase(DB *db.Store) *Database {
	return &Database{DB}
}

func (d *Database) create(ctx context.Context, t *table) (*table, error) {

	type CreateTableResult struct {
		ID        uuid.UUID
		Name      string
		AccountId uuid.NullUUID
	}

	type CreateItemResult struct {
		ID          uuid.UUID
		Description string
		Status      int
	}

	statement := "insert into tables(name, account_id) values($1, $2) returning id, name, account_id"

	resultTable := &CreateTableResult{}
	items := make([]item, len(t.items))

	err := d.DB.ExecTx(ctx, func(tx *sql.Tx) error {

		err := tx.QueryRowContext(ctx, statement, t.name, t.owner.id).Scan(&resultTable.ID, &resultTable.Name, &resultTable.AccountId)

		if err != nil {
			return fmt.Errorf("not was possible to insert the account %v", err)
		}

		itemResult := &CreateItemResult{}

		statement = "insert into items(description, status, table_id) values($1, $2, $3) returning id, description, status"

		for i, current := range t.items {
			err := tx.QueryRowContext(ctx, statement, current.description, current.status, resultTable.ID).Scan(&itemResult.ID, &itemResult.Description, &itemResult.Status)

			if err != nil {
				return fmt.Errorf("not was possible to insert item %v", err)
			}

			items[i] = item{
				id:          itemResult.ID,
				description: itemResult.Description,
				status:      Status(itemResult.Status),
			}
		}

		return nil

	})

	if err != nil {
		return nil, fmt.Errorf("not was possible to create table %v", err)
	}

	return &table{
		id:    resultTable.ID,
		name:  resultTable.Name,
		owner: &owner{resultTable.AccountId},
		items: items,
	}, nil
}

func (d *Database) findByItemId(ctx context.Context, tableID uuid.UUID, id uuid.UUID) (*item, error) {

	type Result struct {
		ID          uuid.UUID
		Description string
		Status      int
		AccountID   uuid.NullUUID
	}

	statement := "select id, description, status, account_id from items where id = $1 and table_id = $2"

	result := &Result{}
	err := d.DB.QueryRowContext(ctx, statement, id, tableID).Scan(&result.ID, &result.Description, &result.Status, &result.AccountID)

	if err != nil {
		return nil, fmt.Errorf("not was possible to find the item %v", err)
	}

	return &item{
		id:          result.ID,
		description: result.Description,
		status:      Status(result.Status),
		owner:       &owner{result.AccountID},
	}, nil
}

func (d *Database) findTableOwnerById(ctx context.Context, tableID uuid.UUID) (*owner, error) {

	type Result struct {
		ID uuid.NullUUID
	}

	statement := "select account_id from tables where id = $1"

	result := &Result{}
	err := d.DB.QueryRowContext(ctx, statement, tableID).Scan(&result.ID)

	if err != nil {
		return nil, fmt.Errorf("not was possible to find the owner %v", err)
	}

	return &owner{result.ID}, nil
}

func (d *Database) updateItem(ctx context.Context, i *item) error {
	return d.DB.ExecTx(ctx, func(t *sql.Tx) error {
		var id uuid.UUID
		err := t.QueryRowContext(ctx, "select id from items where id = $1 for update", i.id).Scan(&id)

		if err != nil {
			return err
		}

		statement := "update items set description=$2, status=$3, account_id=$4 where id = $1"

		_, err = t.ExecContext(ctx, statement, i.id, i.description, i.status, i.owner.id)

		if err != nil {
			return err
		}

		return nil
	})
}
