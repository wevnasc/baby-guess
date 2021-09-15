package tables

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Status int

const (
	None     Status = 0
	Selected        = 1
	Approved        = 2
)

type owner struct {
	id uuid.NullUUID
}

func newOwner(id uuid.UUID) *owner {
	return &owner{uuid.NullUUID{UUID: id, Valid: true}}
}

func (o *owner) isEquals(other *owner) bool {
	return o.id == other.id
}

func (o *owner) nullableID() *uuid.UUID {
	if o.id.Valid {
		return &o.id.UUID
	}

	return nil
}

type item struct {
	id          uuid.UUID
	description string
	luckNumber  int
	owner       *owner
	status      Status
}

func (i *item) selectedBy(owner owner) error {
	if i.status != None {
		return errors.New("just empty items can be selected")
	}
	i.status = Selected
	i.owner = &owner
	return nil
}

func (i *item) unselect() {
	i.status = None
	i.owner = &owner{}
}

func (i *item) approve() error {
	if i.status != Selected {
		return errors.New("just selected items can be approved")
	}
	i.status = Approved
	return nil
}

func (i *item) isOwner(owner *owner) bool {
	return i.owner.isEquals(owner)
}

type table struct {
	id    uuid.UUID
	owner *owner
	name  string
	items []item
}

func newTable(id uuid.UUID, name string, numberItems int) *table {

	items := make([]item, numberItems)

	for i := 0; i < numberItems; i++ {
		key := i + 1
		desc := fmt.Sprintf("%d", key)

		items[i] = item{
			description: desc,
			status:      None,
			luckNumber:  key,
		}
	}

	return &table{name: name, owner: newOwner(id), items: items}
}
