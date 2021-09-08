package tables

import (
	"fmt"

	"github.com/google/uuid"
)

type Status int

const (
	None     Status = 0
	Pending         = 1
	Approved        = 2
)

type item struct {
	id          uuid.UUID
	description string
	ownerID     uuid.UUID
	status      Status
}

type table struct {
	id      uuid.UUID
	ownerID uuid.UUID
	name    string
	items   []item
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

	return &table{name: name, ownerID: id, items: items}
}
