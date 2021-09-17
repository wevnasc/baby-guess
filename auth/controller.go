package auth

import (
	"context"

	"github.com/wevnasc/baby-guess/email"
	"github.com/wevnasc/baby-guess/server"
	"github.com/wevnasc/baby-guess/token"
)

type controller struct {
	database *Database
	email    *email.Connection
}

func newController(database *Database, email *email.Connection) *controller {
	return &controller{database, email}
}

func (c *controller) create(ctx context.Context, account *account) (*account, error) {
	current, _ := c.database.findByEmail(ctx, account.email)

	if current != nil {
		return nil, server.NewError("not was possible to create the account", server.ResourceInvalid)
	}

	a, err := c.database.create(ctx, account)

	if err != nil {
		return nil, server.NewError("not was possible to create the account", server.OperationError)
	}

	go func() {
		e, _ := email.NewFromTemplate(c.email, email.AccountCreated)
		e.Send([]string{a.email}, map[string]string{"name": a.name})
	}()

	return a, nil
}

func (c *controller) login(ctx context.Context, credentials *token.Credentials) (*account, error) {
	account, err := c.database.findByEmail(ctx, credentials.Email)

	if err != nil {
		return nil, server.NewError("not was possible to authenticate the account", server.ResourceInvalid)
	}

	if !account.comparerPassword(credentials.Password) {
		return nil, server.NewError("not was possible to authenticate the account", server.ResourceInvalid)
	}

	return account, nil
}
