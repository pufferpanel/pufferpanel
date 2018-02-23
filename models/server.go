package models

import (
	"github.com/satori/go.uuid"
	"time"
)

type Server struct {
	ID        uuid.UUID `db:"id"`
	Node      Node      `db:node_id belongs_to:node`
	Name      string    `db:name`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Servers []Server
