package tables

import (
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
	id uuid.UUID
}

func (o *owner) isOwner(other *owner) bool {
	return o.id == other.id
}

type item struct {
	id          uuid.UUID
	description string
	owner       *owner
	status      Status
}

func (i *item) selectedBy(owner owner) {
	i.status = Selected
	i.owner = &owner
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
		desc := fmt.Sprintf("%d", i+1)

		items[i] = item{
			description: desc,
			status:      None,
		}
	}

	return &table{name: name, owner: &owner{id}, items: items}
}
