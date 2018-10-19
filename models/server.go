package models

import (
	"github.com/satori/go.uuid"
	"time"
)

type Server struct {
	ID   uint   `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	Name string `gorm:"UNIQUE;size:20"`
	UUID uuid.UUID

	NodeID uint
	Node   Node `gorm:"association_autoupdate:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
