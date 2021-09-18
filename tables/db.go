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

	type AccountResult struct {
		ID    uuid.NullUUID
		Name  string
		Email string
	}

	type TableResult struct {
		ID      uuid.UUID
		Name    string
		account *AccountResult
	}

	tables := []table{}

	statement := "select t.id, t.name, t.account_id, a.name, a.email from tables as t left join accounts as a on t.account_id = a.id where t.account_id = $1"

	tt := []TableResult{}

	rows, err := d.DB.QueryContext(ctx, statement, ownerID)

	if err != nil {
		return nil, fmt.Errorf("not was possible to query tables %v", err)
	}

	for rows.Next() {
		t := &TableResult{account: &AccountResult{}}

		if err := rows.Scan(&t.ID, &t.Name, &t.account.ID, &t.account.Name, &t.account.Email); err != nil {
			return nil, err
		}

		tt = append(tt, *t)
	}

	for _, t := range tt {

		items, err := d.findAllItemsByTableID(ctx, t.ID)

		if err != nil {
			return nil, err
		}

		currentOwner := &owner{
			id:    t.account.ID,
			name:  t.account.Name,
			email: t.account.Email,
		}

		current := table{
			id:    t.ID,
			name:  t.Name,
			owner: currentOwner,
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

	type AccountResult struct {
		ID    uuid.NullUUID
		Name  db.NullString
		Email db.NullString
	}

	type ItemResult struct {
		ID          uuid.UUID
		Description string
		LuckNumber  int
		Winner      bool
		Status      int
		account     *AccountResult
	}

	items := []item{}

	statement := "select i.id, i.description, i.status, i.luck_number, i.winner, i.account_id, a.name, a.email from items as i left join accounts as a on i.account_id = a.id where i.table_id = $1 order by i.luck_number"

	ii := []ItemResult{}

	rows, err := d.DB.QueryContext(ctx, statement, tableID)

	if err != nil {
		return nil, fmt.Errorf("not was possible to query items %v", err)
	}

	for rows.Next() {
		i := ItemResult{account: &AccountResult{}}

		if err := rows.Scan(&i.ID, &i.Description, &i.Status, &i.LuckNumber, &i.Winner, &i.account.ID, &i.account.Name, &i.account.Email); err != nil {
			return nil, err
		}

		ii = append(ii, i)
	}

	for _, i := range ii {
		currentOwner := &owner{
			id:    i.account.ID,
			name:  string(i.account.Name),
			email: string(i.account.Email),
		}

		current := item{
			id:          i.ID,
			description: i.Description,
			status:      Status(i.Status),
			luckNumber:  i.LuckNumber,
			winner:      i.Winner,
			owner:       currentOwner,
		}
		items = append(items, current)
	}

	if err != nil {
		return nil, err
	}

	return items, nil

}

func (d *Database) findByID(ctx context.Context, tableID uuid.UUID) (*table, error) {

	type AccountResult struct {
		ID    uuid.NullUUID
		Name  string
		Email string
	}

	type TableResult struct {
		ID      uuid.UUID
		Name    string
		account *AccountResult
	}

	statement := "select t.id, t.name, t.account_id, a.Name, a.Email from tables as t left join accounts as a on t.account_id = a.id where t.id = $1"

	tr := TableResult{account: &AccountResult{}}

	if err := d.DB.QueryRowContext(ctx, statement, tableID).Scan(&tr.ID, &tr.Name, &tr.account.ID, &tr.account.Name, &tr.account.Email); err != nil {
		return nil, fmt.Errorf("not was possible to query tables %v", err)
	}

	items, err := d.findAllItemsByTableID(ctx, tableID)

	if err != nil {
		return nil, err
	}

	owner := &owner{
		id:    tr.account.ID,
		name:  tr.account.Name,
		email: tr.account.Email,
	}

	return &table{
		id:    tr.ID,
		name:  tr.Name,
		owner: owner,
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

	type AccountResult struct {
		ID    uuid.NullUUID
		Name  db.NullString
		Email db.NullString
	}

	type ItemResult struct {
		ID          uuid.UUID
		Description string
		Status      int
		LuckNumber  int
		Winner      bool
		account     *AccountResult
	}

	statement := "select i.id, i.description, i.status, i.luck_number, i.winner, i.account_id, a.name, a.email from items as i left join accounts as a on i.account_id = a.id where i.id = $1 and i.table_id = $2"

	result := &ItemResult{account: &AccountResult{}}
	err := d.DB.QueryRowContext(ctx, statement, id, tableID).Scan(
		&result.ID,
		&result.Description,
		&result.Status,
		&result.LuckNumber,
		&result.Winner,
		&result.account.ID,
		&result.account.Name,
		&result.account.Email,
	)

	if err != nil {
		return nil, fmt.Errorf("not was possible to find the item %v", err)
	}

	owner := &owner{
		id:    result.account.ID,
		name:  string(result.account.Name),
		email: string(result.account.Email),
	}

	return &item{
		id:          result.ID,
		description: result.Description,
		luckNumber:  result.LuckNumber,
		winner:      result.Winner,
		status:      Status(result.Status),
		owner:       owner,
	}, nil
}

func (d *Database) findTableOwnerByID(ctx context.Context, tableID uuid.UUID) (*owner, error) {

	type Result struct {
		ID    uuid.NullUUID
		Name  string
		Email string
	}

	statement := "select a.id, a.name, a.email from tables as t inner join accounts as a on t.account_id = a.id where t.id = $1"

	result := &Result{}
	err := d.DB.QueryRowContext(ctx, statement, tableID).Scan(&result.ID, &result.Name, &result.Email)

	if err != nil {
		return nil, fmt.Errorf("not was possible to find the owner %v", err)
	}

	return &owner{
		id:    result.ID,
		name:  result.Name,
		email: result.Email,
	}, nil
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

func (d *Database) findOwnerByID(ctx context.Context, ownerID uuid.UUID) (*owner, error) {

	type Result struct {
		ID    uuid.NullUUID
		Name  string
		Email string
	}

	statement := "select a.id, a.name, a.email from accounts as a where a.id = $1"

	result := &Result{}
	err := d.DB.QueryRowContext(ctx, statement, ownerID).Scan(&result.ID, &result.Name, &result.Email)

	if err != nil {
		return nil, fmt.Errorf("not was possible to find the owner %v", err)
	}

	return &owner{
		id:    result.ID,
		name:  result.Name,
		email: result.Email,
	}, nil
}
