package models

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"time"
)

type Server struct {
	Name       string `gorm:"size:40;NOT NULL" json:"-" validate:"required,printascii"`
	Identifier string `gorm:"UNIQUE;NOT NULL;primaryKey;size:8" json:"-" validate:"required,printascii"`

	RawNodeID *uint `gorm:"column:node_id" json:"-"`
	NodeID    uint  `gorm:"-" json:"-" validate:"-"`
	Node      Node  `gorm:"ASSOCIATION_SAVE_REFERENCE:false;foreignKey:RawNodeID" json:"-" validate:"-"`

	IP   string `gorm:"" json:"-" validate:"omitempty,ip|fqdn"`
	Port uint16 `gorm:"" json:"-" validate:"omitempty"`

	Type string `gorm:"NOT NULL;default='generic'" json:"-" validate:"required,printascii"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (s *Server) IsValid() (err error) {
	err = validator.New().Struct(s)
	if err != nil {
		err = pufferpanel.GenerateValidationMessage(err)
	}
	return
}

func (s *Server) BeforeSave(*gorm.DB) (err error) {
	err = s.IsValid()
	if s.NodeID == 0 {
		s.RawNodeID = nil
	} else {
		s.RawNodeID = &s.NodeID
	}
	return
}

func (s *Server) AfterFind(*gorm.DB) (err error) {
	if s.RawNodeID == nil || s.NodeID == LocalNode.ID {
		s.Node = *LocalNode
	}
	return
}
