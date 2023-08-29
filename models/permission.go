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
	"gorm.io/gorm"
	"strings"
)

type Permissions struct {
	ID uint `gorm:"primaryKey,autoIncrement" json:"-"`

	//owners of this permission set
	UserId *uint `json:"-"`
	User   User  `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

	ClientId *uint  `json:"-"`
	Client   Client `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

	//if this set is for a server, what server
	ServerIdentifier *string `json:"-"`
	Server           Server  `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

	RawScopes string              `gorm:"column:scopes;NOT NULL" json:"-" validate:"required"`
	Scopes    []pufferpanel.Scope `gorm:"-" json:"-"`
}

func (p *Permissions) BeforeSave(*gorm.DB) error {
	if p.ServerIdentifier != nil && *p.ServerIdentifier == "" {
		p.ServerIdentifier = nil
	}

	tmp := make([]string, len(p.Scopes))
	for k, v := range p.Scopes {
		tmp[k] = v.String()
	}
	p.RawScopes = strings.Join(tmp, ",")
	return nil
}

func (p *Permissions) AfterFind(*gorm.DB) error {
	p.Scopes = make([]pufferpanel.Scope, 0)
	for _, v := range strings.Split(p.RawScopes, ",") {
		//we can just simply blindly assign it, because the checks we do are just making these strings anyways...
		p.Scopes = append(p.Scopes, pufferpanel.GetScope(v))
	}

	return nil
}

func (p *Permissions) ShouldDelete() bool {
	if p.ServerIdentifier == nil {
		return false
	}
	return len(p.Scopes) == 0
}
