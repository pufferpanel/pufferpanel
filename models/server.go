package models

import (
	"github.com/pufferpanel/pufferpanel/v2"
	"gopkg.in/go-playground/validator.v9"
	"time"
)

type Server struct {
	Name       string `gorm:"UNIQUE_INDEX;size:20;NOT NULL" json:"-" validate:"required,printascii"`
	Identifier string `gorm:"UNIQUE_INDEX;NOT NULL;PRIMARY_KEY;size:8" json:"-" validate:"required,printascii"`

	NodeID uint `gorm:"NOT NULL" json:"-" validate:"required,min=1"`
	Node   Node `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

	IP   string `gorm:"" json:"-" validate:"required,ip|fqdn"`
	Port uint   `gorm:"" json:"-" validate:"required,min=1,max=65535"`

	Type string `gorm:"NOT NULL;default='generic'" json:"-" validate:"required,printascii"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Servers []*Server

func (s *Server) IsValid() (err error) {
	err = validator.New().Struct(s)
	if err != nil {
		err = pufferpanel.GenerateValidationMessage(err)
	}
	return
}

func (s *Server) BeforeSave() (err error) {
	err = s.IsValid()
	return
}
