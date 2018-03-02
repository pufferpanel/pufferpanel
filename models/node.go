package models

import (
	"github.com/gobuffalo/uuid"
	"time"
)

type Node struct {
	ID        uuid.UUID `db:"id"`
	Location  Location  `db:location_id belongs_to:location`
	Name      string    `db:name`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Nodes []Node

func CreateNode(location Location, name string) (node Node, err error) {
	id, err := uuid.NewV4()
	if err != nil {
		return
	}

	node = Node{
		ID:       id,
		Location: location,
		Name:     name,
	}

	return
}

func (n *Node) Delete() (err error) {
	err = DB.Destroy(n)
	return
}

func (n *Node) Save() (err error) {
	err = DB.Save(n)
	return
}
