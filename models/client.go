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
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"strings"
)

type Client struct {
	ID                 uint   `gorm:"PRIMARY_KEY,AUTO_INCREMEMT" json:"-"`
	ClientId           string `gorm:"NOT NULL" json:"client_id"`
	HashedClientSecret string `gorm:"column:client_secret;NOT NULL" json:"-"`

	UserId uint  `gorm:"NOT NULL" json:"-"`
	User   *User `json:"-"`

	ServerId string  `gorm:"NOT NULL" json:"-"`
	Server   *Server `json:"-"`

	Scopes    []pufferpanel.Scope `gorm:"-" json:"-"`
	RawScopes string              `gorm:"column:scopes;NOT NULL;size:4000" json:"-"`

	Name        string `gorm:"column:name;NOT NULL;size:100;default\"\"" json:"name"`
	Description string `gorm:"column:description;NOT NULL;size:4000;default:\"\"" json:"description"`
}

type CreatedClient struct {
	ClientId     string `json:"id"`
	ClientSecret string `json:"secret"`
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

func (c *Client) BeforeSave(*gorm.DB) (err error) {
	err = c.IsValid()

	scopes := make([]string, 0)

	for _, s := range c.Scopes {
		scopes = append(scopes, string(s))
	}
	c.RawScopes = strings.Join(scopes, " ")

	return
}

func (c *Client) AfterFind(*gorm.DB) (err error) {
	if c.RawScopes == "" {
		c.Scopes = make([]pufferpanel.Scope, 0)
		return
	}

	split := strings.Split(c.RawScopes, " ")
	c.Scopes = make([]pufferpanel.Scope, len(split))

	for i, v := range split {
		c.Scopes[i] = pufferpanel.Scope(v)
	}

	return
}
