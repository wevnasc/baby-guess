package tables

import "context"

type controller struct {
	database *Database
}

func newController(database *Database) *controller {
	return &controller{database}
}

func (c *controller) create(ctx context.Context, table *table) (*table, error) {
	return c.database.create(ctx, table)
}
