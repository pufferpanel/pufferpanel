package models

import (
	"github.com/satori/go.uuid"
	"time"
)

type Server struct {
	ID        uuid.UUID `db:"id"`
	NodeID    uuid.UUID `db:node_id`
	Name      string    `db:name`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Servers []Server
