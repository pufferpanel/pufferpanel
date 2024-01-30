package models

import (
	uuid "github.com/gofrs/uuid/v5"
	"github.com/pufferpanel/pufferpanel/v3"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Node struct {
	ID          uint   `json:"-"`
	Name        string `gorm:"size:100;UNIQUE;NOT NULL" json:"-" validate:"required,printascii"`
	PublicHost  string `gorm:"size:100;NOT NULL" json:"-" validate:"required,ip|fqdn"`
	PrivateHost string `gorm:"size:100;NOT NULL" json:"-" validate:"required,ip|fqdn"`
	PublicPort  uint16 `gorm:"DEFAULT:8008;NOT NULL" json:"-" validate:"required,min=1,max=65535,nefield=SFTPPort"`
	PrivatePort uint16 `gorm:"DEFAULT:8008;NOT NULL" json:"-" validate:"required,min=1,max=65535,nefield=SFTPPort"`
	SFTPPort    uint16 `gorm:"DEFAULT:5657;NOT NULL" json:"-" validate:"required,min=1,max=65535,nefield=PublicPort,nefield=PrivatePort"`

	Secret string `gorm:"size=36;NOT NULL" json:"-" validate:"required"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (n *Node) IsValid() (err error) {
	err = validator.New().Struct(n)
	if err != nil {
		err = pufferpanel.GenerateValidationMessage(err)
	}
	return
}

func (n *Node) BeforeSave(*gorm.DB) (err error) {
	err = n.IsValid()
	return
}

func (n *Node) IsLocal() bool {
	return n.ID == LocalNode.ID
}

var LocalNode = &Node{
	ID:          0,
	Name:        "LocalNode",
	PublicHost:  "127.0.0.1",
	PrivateHost: "127.0.0.1",
	PublicPort:  8008,
	PrivatePort: 8008,
	SFTPPort:    5657,
}

func init() {
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	LocalNode.Secret = strings.Replace(u.String(), "-", "", -1)
}
