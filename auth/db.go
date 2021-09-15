package auth

import (
	"context"
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

func (d *Database) create(ctx context.Context, a *account) (*account, error) {
	type CreateResult struct {
		ID    string
		Name  string
		Email string
	}

	statement := "insert into accounts(name, email, password) values($1, $2, $3) returning id, name, email"

	result := &CreateResult{}
	err := d.DB.QueryRowContext(ctx, statement, a.name, a.email, a.password).Scan(&result.ID, &result.Name, &result.Email)

	if err != nil {
		return nil, fmt.Errorf("not was possible to insert the account %v", err)
	}

	uuid, err := uuid.Parse(result.ID)

	if err != nil {
		return nil, fmt.Errorf("error to generate uuid %v", err)
	}

	return &account{
		id:    uuid,
		name:  result.Name,
		email: result.Email,
	}, nil
}

func (d *Database) findByEmail(ctx context.Context, email string) (*account, error) {
	type FindByEmailResult struct {
		ID    string
		Name  string
		Email string
	}

	statement := "select id, name, email from accounts where email = $1"

	result := &FindByEmailResult{}
	err := d.DB.QueryRowContext(ctx, statement, email).Scan(&result.ID, &result.Name, &result.Email)

	if err != nil {
		return nil, fmt.Errorf("not was possible to find the account %v", err)
	}

	uuid, err := uuid.Parse(result.ID)

	if err != nil {
		return nil, fmt.Errorf("error to generate uuid %v", err)
	}

	return &account{
		id:    uuid,
		name:  result.Name,
		email: result.Email,
	}, nil
}
