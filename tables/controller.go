package tables

import (
	"context"

	"github.com/google/uuid"
	"github.com/wevnasc/baby-guess/email"
	"github.com/wevnasc/baby-guess/server"
)

type controller struct {
	database *Database
	email    email.Client
}

func newController(database *Database, email email.Client) *controller {
	return &controller{database, email}
}

func (c *controller) create(ctx context.Context, table *table) (*table, error) {
	return c.database.create(ctx, table)
}

func (c *controller) all(ctx context.Context, accountID uuid.UUID) ([]table, error) {
	return c.database.findAllByOwnerID(ctx, accountID)
}

func (c *controller) selectItem(ctx context.Context, tableID uuid.UUID, selected item) error {
	owner, err := c.database.findTableOwnerByID(ctx, tableID)

	if err != nil {
		return server.NewError("table not found", server.ResourceNotFound)
	}

	if owner.isEquals(selected.owner) {
		return server.NewError("the table's owner can't select an item", server.OperationNotAllowed)
	}

	item, err := c.database.findByItemID(ctx, tableID, selected.id)

	if err != nil {
		return server.NewError("item not found", server.ResourceNotFound)
	}

	if err := item.selectedBy(*selected.owner); err != nil {
		return server.NewError(err.Error(), server.OperationNotAllowed)
	}

	if err := c.database.updateItem(ctx, item); err != nil {
		return server.NewError("not was possible to select the item", server.OperationError)
	}

	go func() {
		c.email.Send(email.ItemSelected, []string{owner.email}, map[string]string{"item": item.description})
	}()

	return nil
}

func (c *controller) unselectItem(ctx context.Context, tableID uuid.UUID, unselected item) error {

	owner, err := c.database.findTableOwnerByID(ctx, tableID)

	if err != nil {
		return server.NewError("table not found", server.ResourceNotFound)
	}

	item, err := c.database.findByItemID(ctx, tableID, unselected.id)

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

	tableOwner, err := c.database.findTableOwnerByID(ctx, tableID)

	if err != nil {
		return server.NewError("table not found", server.ResourceNotFound)
	}

	item, err := c.database.findByItemID(ctx, tableID, itemID)

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

func (c *controller) draw(ctx context.Context, owner *owner, tableID uuid.UUID) (*item, error) {

	table, err := c.database.findByID(ctx, tableID)

	if err != nil {
		return nil, server.NewError("table not found", server.ResourceNotFound)
	}

	if !table.isOwner(owner) {
		return nil, server.NewError("just the table's owner can draw a winner", server.OperationNotAllowed)
	}

	item, err := table.drawWinner()

	if err != nil {
		return nil, server.NewError(err.Error(), server.OperationError)
	}

	if err := c.database.updateItem(ctx, item); err != nil {
		return nil, server.NewError(err.Error(), server.OperationError)
	}

	return item, nil
}
