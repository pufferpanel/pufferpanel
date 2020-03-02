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
	"github.com/pufferpanel/pufferpanel/v2"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"time"
)

type User struct {
	ID             uint   `json:"-"`
	Username       string `gorm:"UNIQUE_INDEX;NOT NULL;size:100" json:"-" validate:"required,printascii,max=100,min=5"`
	Email          string `gorm:"UNIQUE_INDEX;NOT NULL;size:255" json:"-" validate:"required,email,max=255"`
	HashedPassword string `gorm:"column:password;NOT NULL;size:200" json:"-" validate:"required,max=200"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Users []*User

func (u *User) SetPassword(pw string) error {
	res, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)

	if err == nil {
		u.HashedPassword = string(res)
	}

	return err
}

func (u *User) IsValid() (err error) {
	err = validator.New().Struct(u)
	if err != nil {
		err = pufferpanel.GenerateValidationMessage(err)
	}
	return
}

func (u *User) BeforeSave() (err error) {
	err = u.IsValid()
	return
}

