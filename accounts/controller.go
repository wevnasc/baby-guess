package accounts

import (
	"context"
	"errors"
)

type controller struct {
	database *Database
}

func newController(database *Database) *controller {
	return &controller{database}
}

func (c *controller) create(ctx context.Context, account *account) (*account, error) {
	current, _ := c.database.findByEmail(ctx, account.email)

	if current == nil {
		return c.database.create(ctx, account)
	}

	return nil, errors.New("not was possible to create the account")
}
