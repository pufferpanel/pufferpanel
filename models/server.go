/*
 Copyright 2020 Padduck, LLC
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
	"github.com/pufferpanel/pufferpanel/v3"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"time"
)

type Server struct {
	Name       string `gorm:"size:40;NOT NULL" json:"-" validate:"required,printascii"`
	Identifier string `gorm:"UNIQUE;NOT NULL;primaryKey;size:8" json:"-" validate:"required,printascii"`

	NodeID uint `gorm:"NOT NULL;default=0" json:"-" validate:"min=0"`
	Node   Node `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

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
	return
}

func (s *Server) AfterFind(*gorm.DB) (err error) {
	if s.NodeID == LocalNode.ID {
		s.Node = *LocalNode
	}
	return
}
