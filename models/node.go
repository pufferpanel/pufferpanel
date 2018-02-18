package models

import (
	"github.com/satori/go.uuid"
	"time"
)

type Node struct {
	ID        uuid.UUID `db:"id"`
	LocationID    uuid.UUID `db:location_id`
	Name      string    `db:name`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Nodes []Node
