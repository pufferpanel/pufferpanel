package models

import (
	"github.com/pufferpanel/pufferpanel/shared"
	"gopkg.in/go-playground/validator.v9"
	"time"
)

type Server struct {
	ID         uint   `json:"-"`
	Name       string `gorm:"UNIQUE_INDEX;size:20;NOT NULL" json:"-"`
	Identifier string `gorm:"UNIQUE_INDEX;NOT NULL;size:8" json:"-"`

	NodeID uint `gorm:"NOT NULL" json:"-"`
	Node   Node `gorm:"association_autoupdate:false" json:"-"`

	IP   string `gorm:"" json:"-"`
	Port uint   `gorm:"" json:"-"`

	Type string `gorm:"NOT NULL;default='generic'" json:"-"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Servers []*Server

func (s *Server) IsValid() (err error) {
	err = validator.New().Struct(s)
	if err != nil {
		err = shared.GenerateValidationMessage(err)
	}
	return
}

func (s *Server) BeforeSave() (err error) {
	err = s.IsValid()
	return
}
