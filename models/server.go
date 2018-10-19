package models

import (
	"github.com/satori/go.uuid"
	"time"
)

type Server struct {
	Id    int    `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	Name  string `gorm:"UNIQUE;size:20"`
	UUID  uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time

	NodeId int
	Node   Node
}
