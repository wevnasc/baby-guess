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

func (d *Database) findAllByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]table, error) {

	type TableResult struct {
		ID        uuid.UUID
		Name      string
		AccountID uuid.NullUUID
	}

	tables := []table{}

	statement := "select id, name, account_id from tables where account_id = $1"

	tt := []TableResult{}

	rows, err := d.DB.QueryContext(ctx, statement, ownerID)

	if err != nil {
		return nil, fmt.Errorf("not was possible to query tables %v", err)
	}

	for rows.Next() {
		t := &TableResult{}

		if err := rows.Scan(&t.ID, &t.Name, &t.AccountID); err != nil {
			return nil, err
		}

		tt = append(tt, *t)
	}

	for _, t := range tt {

		items, err := d.findAllItemsByTableID(ctx, t.ID)

		if err != nil {
			return nil, err
		}

		current := table{
			id:    t.ID,
			name:  t.Name,
			owner: &owner{id: t.AccountID},
			items: items,
		}

		tables = append(tables, current)
	}

	if err != nil {
		return nil, err
	}

	return tables, nil

}

func (d *Database) findAllItemsByTableID(ctx context.Context, tableID uuid.UUID) ([]item, error) {

	type ItemResult struct {
		ID          uuid.UUID
		Description string
		LuckNumber  int
		Winner      bool
		Status      int
		AccountID   uuid.NullUUID
	}

	items := []item{}

	statement := "select id, description, status, luck_number, winner, account_id from items where table_id = $1 order by luck_number"

	ii := []ItemResult{}

	rows, err := d.DB.QueryContext(ctx, statement, tableID)

	if err != nil {
		return nil, fmt.Errorf("not was possible to query items %v", err)
	}

	for rows.Next() {
		i := ItemResult{}

		if err := rows.Scan(&i.ID, &i.Description, &i.Status, &i.LuckNumber, &i.Winner, &i.AccountID); err != nil {
			return nil, err
		}

		ii = append(ii, i)
	}

	for _, i := range ii {
		current := item{
			id:          i.ID,
			description: i.Description,
			status:      Status(i.Status),
			luckNumber:  i.LuckNumber,
			winner:      i.Winner,
			owner:       &owner{id: i.AccountID},
		}
		items = append(items, current)
	}

	if err != nil {
		return nil, err
	}

	return items, nil

}

func (d *Database) findByID(ctx context.Context, tableID uuid.UUID) (*table, error) {
	type Result struct {
		ID        uuid.UUID
		Name      string
		AccountID uuid.NullUUID
	}

	statement := "select id, name, account_id from tables where id = $1"

	result := Result{}

	if err := d.DB.QueryRowContext(ctx, statement, tableID).Scan(&result.ID, &result.Name, &result.AccountID); err != nil {
		return nil, fmt.Errorf("not was possible to query tables %v", err)
	}

	items, err := d.findAllItemsByTableID(ctx, tableID)

	if err != nil {
		return nil, err
	}

	return &table{
		id:    result.ID,
		name:  result.Name,
		owner: &owner{id: result.AccountID},
		items: items,
	}, nil
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
		LuckNumber  int
		Winner      bool
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

		statement = "insert into items(description, status, luck_number, table_id) values($1, $2, $3, $4) returning id, description, status, luck_number, winner"

		for i, current := range t.items {
			err := tx.QueryRowContext(ctx, statement, current.description, current.status, current.luckNumber, resultTable.ID).Scan(&itemResult.ID, &itemResult.Description, &itemResult.Status, &itemResult.LuckNumber, &itemResult.Winner)

			if err != nil {
				return fmt.Errorf("not was possible to insert item %v", err)
			}

			items[i] = item{
				id:          itemResult.ID,
				description: itemResult.Description,
				luckNumber:  itemResult.LuckNumber,
				winner:      itemResult.Winner,
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
		owner: &owner{id: resultTable.AccountId},
		items: items,
	}, nil
}

func (d *Database) findByItemID(ctx context.Context, tableID uuid.UUID, id uuid.UUID) (*item, error) {

	type Result struct {
		ID          uuid.UUID
		Description string
		Status      int
		LuckNumber  int
		Winner      bool
		AccountID   uuid.NullUUID
	}

	statement := "select id, description, status, luck_number, winner, account_id from items where id = $1 and table_id = $2"

	result := &Result{}
	err := d.DB.QueryRowContext(ctx, statement, id, tableID).Scan(&result.ID, &result.Description, &result.Status, &result.LuckNumber, &result.Winner, &result.AccountID)

	if err != nil {
		return nil, fmt.Errorf("not was possible to find the item %v", err)
	}

	return &item{
		id:          result.ID,
		description: result.Description,
		luckNumber:  result.LuckNumber,
		winner:      result.Winner,
		status:      Status(result.Status),
		owner:       &owner{id: result.AccountID},
	}, nil
}

func (d *Database) findTableOwnerByID(ctx context.Context, tableID uuid.UUID) (*owner, error) {

	type Result struct {
		ID    uuid.NullUUID
		Email string
	}

	statement := "select a.id, a.email from tables as t inner join accounts as a on t.account_id = a.id where t.id = $1"

	result := &Result{}
	err := d.DB.QueryRowContext(ctx, statement, tableID).Scan(&result.ID, &result.Email)

	if err != nil {
		return nil, fmt.Errorf("not was possible to find the owner %v", err)
	}

	return &owner{result.ID, result.Email}, nil
}

func (d *Database) updateItem(ctx context.Context, i *item) error {
	return d.DB.ExecTx(ctx, func(t *sql.Tx) error {
		var id uuid.UUID
		err := t.QueryRowContext(ctx, "select id from items where id = $1 for update", i.id).Scan(&id)

		if err != nil {
			return err
		}

		statement := "update items set description=$2, status=$3, luck_number=$4, winner=$5, account_id=$6 where id = $1"

		_, err = t.ExecContext(ctx, statement, i.id, i.description, i.status, i.luckNumber, i.winner, i.owner.id)

		if err != nil {
			return err
		}

		return nil
	})
}
