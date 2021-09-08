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
		ID        string
		Name      string
		AccountId string
	}

	type CreateItemResult struct {
		ID          string
		Description string
		Status      int
	}

	tableStmt := "insert into tables(name, account_id) values($1, $2) returning id, name, account_id"
	itemStmt := "insert into items(description, status, table_id) values($1, $2, $3) returning id, description, status"

	resultTable := &CreateTableResult{}
	items := make([]item, len(t.items))

	err := d.DB.ExecTx(ctx, func(tx *sql.Tx) error {

		err := tx.QueryRowContext(ctx, tableStmt, t.name, t.ownerID).Scan(&resultTable.ID, &resultTable.Name, &resultTable.AccountId)

		if err != nil {
			return fmt.Errorf("not was possible to insert the account %v", err)
		}

		itemResult := &CreateItemResult{}

		for i, current := range t.items {
			err := tx.QueryRowContext(ctx, itemStmt, current.description, current.status, resultTable.ID).Scan(&itemResult.ID, &itemResult.Description, &itemResult.Status)

			if err != nil {
				return fmt.Errorf("not was possible to insert item %v", err)
			}

			uuid, err := uuid.Parse(itemResult.ID)

			fmt.Println(itemResult.Description)

			if err != nil {
				return fmt.Errorf("error to generate uuid %v", err)
			}

			items[i] = item{
				id:          uuid,
				description: itemResult.Description,
				status:      Status(itemResult.Status),
			}
		}

		return nil

	})

	if err != nil {
		return nil, fmt.Errorf("not was possible to create table %v", err)
	}

	tableUUID, err := uuid.Parse(resultTable.ID)

	if err != nil {
		return nil, fmt.Errorf("error to generate uuid %v", err)
	}

	accountUUID, err := uuid.Parse(resultTable.AccountId)

	if err != nil {
		return nil, fmt.Errorf("error to generate uuid %v", err)
	}

	return &table{
		id:      tableUUID,
		name:    resultTable.Name,
		ownerID: accountUUID,
		items:   items,
	}, nil
}
