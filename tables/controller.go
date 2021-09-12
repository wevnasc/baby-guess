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

	if owner.isEquals(selected.owner) {
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

func (c *controller) unselectItem(ctx context.Context, tableID uuid.UUID, unselected item) error {

	owner, err := c.database.findTableOwnerById(ctx, tableID)

	if err != nil {
		return server.NewError("table not found", server.ResourceNotFound)
	}

	item, err := c.database.findByItemId(ctx, tableID, unselected.id)

	if err != nil {
		return server.NewError("item not found", server.ResourceNotFound)
	}

	if !(owner.isEquals(unselected.owner) || item.isOwner(unselected.owner)) {
		return server.NewError("just the table's owner or the item's owner can unselect the item", server.OperationNotAllowed)
	}

	item.unselect()

	if err := c.database.updateItem(ctx, item); err != nil {
		return server.NewError("not was possible to unselect the item", server.OperationError)
	}

	return nil
}

func (c *controller) approveItem(ctx context.Context, owner *owner, tableID uuid.UUID, itemID uuid.UUID) error {

	tableOwner, err := c.database.findTableOwnerById(ctx, tableID)

	if err != nil {
		return server.NewError("table not found", server.ResourceNotFound)
	}

	item, err := c.database.findByItemId(ctx, tableID, itemID)

	if err != nil {
		return server.NewError("item not found", server.ResourceNotFound)
	}

	if !tableOwner.isEquals(owner) {
		return server.NewError("just the table's owner can approve an item", server.OperationNotAllowed)
	}

	if err := item.approve(); err != nil {
		return server.NewError(err.Error(), server.OperationNotAllowed)
	}

	if err := c.database.updateItem(ctx, item); err != nil {
		return server.NewError("not was possible to unselect the item", server.OperationError)
	}

	return nil
}
