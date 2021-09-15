package auth

import (
	"context"

	"github.com/wevnasc/baby-guess/server"
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

	return nil, server.NewError("not was possible to create the account", server.ResourceInvalid)
}

func (c *controller) login(ctx context.Context, credentials *credentials) (*account, error) {
	account, err := c.database.findByEmail(ctx, credentials.email)

	if err != nil {
		return nil, server.NewError("not was possible to authenticate the account", server.ResourceInvalid)
	}

	if !account.comparerPassword(credentials.password) {
		return nil, server.NewError("not was possible to authenticate the account", server.ResourceInvalid)
	}

	return account, nil
}
