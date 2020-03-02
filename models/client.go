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
	"strings"
)

type Client struct {
	ID                 uint   `gorm:"PRIMARY_KEY,AUTO_INCREMEMT" json:"-"`
	ClientId           string `gorm:"NOT NULL"`
	HashedClientSecret string `gorm:"column:client_secret;NOT NULL"`

	UserId uint `gorm:"NOT NULL"`
	User   *User

	ServerId string `gorm:"NOT NULL"`
	Server   *Server

	Scopes    []pufferpanel.Scope `gorm:"-"`
	RawScopes string        `gorm:"column:scopes;NOT NULL;size:4000"`
}

type Clients []*Client

func (c *Client) SetClientSecret(secret string) error {
	res, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)

	if err == nil {
		c.HashedClientSecret = string(res)
	}

	return err
}

func (c *Client) ValidateSecret(secret string) bool {
	return bcrypt.CompareHashAndPassword([]byte(c.HashedClientSecret), []byte(secret)) == nil
}

func (c *Client) IsValid() (err error) {
	err = validator.New().Struct(c)

	if err != nil {
		err = pufferpanel.GenerateValidationMessage(err)
	}

	return
}

func (c *Client) BeforeSave() (err error) {
	err = c.IsValid()

	scopes := make([]string, 0)

	for _, s := range c.Scopes {
		scopes = append(scopes, string(s))
	}
	c.RawScopes = strings.Join(scopes, " ")

	return
}

func (c *Client) AfterFind() (err error) {
	split := strings.Split(c.RawScopes, " ")
	c.Scopes = make([]pufferpanel.Scope, len(split))

	for i, v := range split {
		c.Scopes[i] = pufferpanel.Scope(v)
	}

	return
}
