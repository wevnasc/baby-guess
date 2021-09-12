package tables

import (
	"context"

	"github.com/google/uuid"
	"github.com/wevnasc/baby-guess/server"
)

type controller struct {
	database *Database
}

func newController(database *Database) *controller {
	return &controller{database}
}

func (c *controller) create(ctx context.Context, table *table) (*table, error) {
	return c.database.create(ctx, table)
}

func (c *controller) selectItem(ctx context.Context, tableID uuid.UUID, selected item) error {

	owner, err := c.database.findTableOwnerById(ctx, tableID)

	if err != nil {
		return server.NewError("table not found", server.ResourceNotFound)
	}

	if owner.isOwner(selected.owner) {
		return server.NewError("the table's owner can't select an item", server.OperationNotAllowed)
	}

	item, err := c.database.findByItemId(ctx, tableID, selected.id)

	if err != nil {
		return server.NewError("item not found", server.ResourceNotFound)
	}

	if err := item.selectedBy(*selected.owner); err != nil {
		return server.NewError(err.Error(), server.OperationNotAllowed)
	}

	if err := c.database.updateItem(ctx, item); err != nil {
		return server.NewError("not was possible to select the item", server.OperationError)
	}

	return nil
}
