package models

import (
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/shared"
	"gopkg.in/go-playground/validator.v9"
)

type Server struct {
	ID         uint   `json:"-"`
	Name       string `gorm:"UNIQUE_INDEX;size:20;NOT NULL" json:"-"`
	Identifier string `gorm:"UNIQUE_INDEX;NOT NULL;size:8" json:"-"`

	NodeID uint `gorm:"NOT NULL" json:"-"`
	Node   Node `gorm:"association_autoupdate:false" json:"-"`

	//CreatedAt time.Time `json:"-"`
	//UpdatedAt time.Time `json:"-"`
}

type Servers []*Server

func MigrateServerModel(db *gorm.DB) (err error) {
	err = db.Model(&Server{}).AddForeignKey("node_id", "nodes(id)", "RESTRICT", "RESTRICT").Error
	return
}

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
