package tables

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Status int

const (
	None     Status = 0
	Selected        = 1
	Approved        = 2
)

const (
	MaxItems = 500
	MinItems = 30
)

type owner struct {
	id    uuid.NullUUID
	name  string
	email string
}

func newOwner(id uuid.UUID) *owner {
	return &owner{id: uuid.NullUUID{UUID: id, Valid: true}}
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
	winner      bool
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

func newTable(id uuid.UUID, name string, numberItems int) (*table, error) {

	if numberItems < MinItems {
		return nil, errors.New("items out of range, the minimum of items should be at least 30")
	}

	if numberItems > MaxItems {
		return nil, errors.New("items out of range, too many items the maximum of items is 500")
	}

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

	return &table{name: name, owner: newOwner(id), items: items}, nil
}

func (t *table) isOwner(other *owner) bool {
	return t.owner.isEquals(other)
}

func (t *table) winner() *owner {
	for _, item := range t.items {
		if item.winner {
			return item.owner
		}
	}

	return nil
}

func (t *table) drawWinner() (*item, error) {
	for _, item := range t.items {
		if item.status != Approved {
			return nil, errors.New("all item should be approved before draw")
		}
	}

	if t.winner() != nil {
		return nil, errors.New("error to draw, winner already exists")
	}

	rand.Seed(time.Now().Unix())
	n := rand.Intn(len(t.items)-1) + 1

	for i, _ := range t.items {
		if t.items[i].luckNumber == n {
			t.items[i].winner = true
			return &t.items[i], nil
		}
	}

	return nil, errors.New("not was possible to draw a number")
}

func (t *table) losers() (list []owner, _ error) {
	exists := make(map[string]bool, len(t.items))

	winner := t.winner()

	if winner == nil {
		return nil, errors.New("no winner found to determine the losers")
	}

	exists[winner.id.UUID.String()] = true

	for _, item := range t.items {

		key := item.owner.id.UUID.String()

		_, ok := exists[key]

		if ok {
			continue
		}
		exists[key] = true

		list = append(list, *item.owner)
	}

	return list, nil
}
