/*
 Copyright 2022 Padduck, LLC
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
 	http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package models

import (
	"github.com/pufferpanel/pufferpanel/v2"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Node struct {
	ID          uint   `json:"-"`
	Name        string `gorm:"size:100;UNIQUE;NOT NULL" json:"-" validate:"required,printascii"`
	PublicHost  string `gorm:"size:100;NOT NULL" json:"-"`
	PrivateHost string `gorm:"size:100;NOT NULL" json:"-"`
	PublicPort  uint16 `gorm:"DEFAULT:8080;NOT NULL" json:"-" validate:"required,min=1,max=65535,nefield=SFTPPort"`
	PrivatePort uint16 `gorm:"DEFAULT:8080;NOT NULL" json:"-" validate:"required,min=1,max=65535,nefield=SFTPPort"`
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
	PublicPort:  8080,
	PrivatePort: 8080,
	SFTPPort:    5657,
	Secret:      strings.Replace(uuid.NewV4().String(), "-", "", -1),
}
