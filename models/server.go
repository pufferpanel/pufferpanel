package models

import (
	"github.com/satori/go.uuid"
	"time"
)

type Server struct {
	ID   uint      `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	Name string    `gorm:"UNIQUE_INDEX;size:20;NOT NULL"`
	UUID uuid.UUID `gorm:"UNIQUE_INDEX;NOT NULL"`

	NodeID uint	`gorm:"NOT NULL"`
	Node   Node `gorm:"association_autoupdate:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Servers []*Server