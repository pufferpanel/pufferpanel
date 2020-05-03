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
	"reflect"
)

type Permissions struct {
	ID uint `gorm:"PRIMARY_KEY,AUTO_INCREMEMT" json:"-"`

	//owners of this permission set
	UserId *uint `json:"-"`
	User   User  `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

	ClientId *uint  `json:"-"`
	Client   Client `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

	//if this set is for a server, what server
	ServerIdentifier *string `json:"-"`
	Server           Server  `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

	//and here are all the perms we support
	Admin           bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	ViewServer      bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	CreateServer    bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	ViewNodes       bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	EditNodes       bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	DeployNodes     bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	ViewTemplates   bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	EditTemplates   bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	EditUsers       bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	ViewUsers       bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	EditServerAdmin bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	DeleteServer    bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`

	//these only will exist if tied to a server, and for a user
	EditServerData    bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	EditServerUsers   bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	InstallServer     bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	UpdateServer      bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""` //this is unused currently
	ViewServerConsole bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	SendServerConsole bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	StopServer        bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	StartServer       bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	ViewServerStats   bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	ViewServerFiles   bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	SFTPServer        bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
	PutServerFiles    bool `gorm:"NOT NULL;DEFAULT:0" json:"-" oneOf:""`
}

type MultiplePermissions []*Permissions

func (p *Permissions) BeforeSave() {
	if p.ServerIdentifier != nil && *p.ServerIdentifier == "" {
		p.ServerIdentifier = nil
	}
}

func (p *Permissions) ToScopes() []pufferpanel.Scope {
	scopes := make([]pufferpanel.Scope, 0)

	if p.Admin {
		scopes = append(scopes, pufferpanel.ScopeServersAdmin)

		if p.ServerIdentifier == nil {
			scopes = append(scopes, pufferpanel.ScopeServersCreate, pufferpanel.ScopeNodesView, pufferpanel.ScopeNodesDeploy, pufferpanel.ScopeNodesEdit, pufferpanel.ScopeTemplatesView, pufferpanel.ScopeUsersView, pufferpanel.ScopeUsersEdit)
		} else {
			scopes = append(scopes, pufferpanel.ScopeServersDelete, pufferpanel.ScopeServersEditAdmin)
		}
	} else {
		if p.ServerIdentifier == nil {
			if p.CreateServer {
				scopes = append(scopes, pufferpanel.ScopeServersCreate)
			}

			if p.ViewNodes {
				scopes = append(scopes, pufferpanel.ScopeNodesView)
			}

			if p.EditNodes {
				scopes = append(scopes, pufferpanel.ScopeNodesEdit)
			}

			if p.ViewTemplates {
				scopes = append(scopes, pufferpanel.ScopeTemplatesView)
			}

			if p.EditTemplates {
				scopes = append(scopes, pufferpanel.ScopeTemplatesEdit)
			}

			if p.EditUsers {
				scopes = append(scopes, pufferpanel.ScopeUsersEdit)
			}

			if p.ViewUsers {
				scopes = append(scopes, pufferpanel.ScopeUsersView)
			}

			if p.DeployNodes {
				scopes = append(scopes, pufferpanel.ScopeNodesDeploy)
			}
		} else {
			if p.DeleteServer {
				scopes = append(scopes, pufferpanel.ScopeServersDelete)
			}

			if p.EditServerAdmin {
				scopes = append(scopes, pufferpanel.ScopeServersEditAdmin)
			}
		}
	}

	if p.ViewServer {
		scopes = append(scopes, pufferpanel.ScopeServersView)
	}

	//these only apply if there is a server involved
	if p.ServerIdentifier != nil {
		if p.EditServerData {
			scopes = append(scopes, pufferpanel.ScopeServersEdit)
		}

		if p.EditServerUsers {
			scopes = append(scopes, pufferpanel.ScopeServersEditUsers)
		}

		if p.InstallServer {
			scopes = append(scopes, pufferpanel.ScopeServersInstall)
		}

		if p.UpdateServer {
			scopes = append(scopes, pufferpanel.ScopeServersUpdate)
		}

		if p.ViewServerConsole {
			scopes = append(scopes, pufferpanel.ScopeServersConsole)
		}

		if p.SendServerConsole {
			scopes = append(scopes, pufferpanel.ScopeServersConsoleSend)
		}

		if p.StopServer {
			scopes = append(scopes, pufferpanel.ScopeServersStop)
		}

		if p.StartServer {
			scopes = append(scopes, pufferpanel.ScopeServersStart)
		}

		if p.ViewServerStats {
			scopes = append(scopes, pufferpanel.ScopeServersStat)
		}

		if p.ViewServerFiles {
			scopes = append(scopes, pufferpanel.ScopeServersFilesGet)
		}

		if p.PutServerFiles {
			scopes = append(scopes, pufferpanel.ScopeServersFilesPut)
		}

		if p.SFTPServer {
			scopes = append(scopes, pufferpanel.ScopeServersSFTP)
		}
	}

	return scopes
}

func (p *Permissions) SetDefaults() {
	p.ViewServer = true

	if p.ServerIdentifier != nil {
		p.EditServerData = true
		p.EditServerUsers = true
		p.InstallServer = true
		p.UpdateServer = true
		p.ViewServerConsole = true
		p.SendServerConsole = true
		p.StopServer = true
		p.StartServer = true
		p.ViewServerStats = true
		p.ViewServerFiles = true
		p.SFTPServer = true
		p.PutServerFiles = true
	}
}

func (p *Permissions) ShouldDelete() bool {
	val := reflect.ValueOf(p)

	// If it's an interface or a pointer, unwrap it.
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldKind := field.Kind()

		if fieldKind != reflect.Bool {
			continue
		}

		typeField := val.Type().Field(i)

		// Get the field tag value.
		_, exist := typeField.Tag.Lookup("oneOf")

		if !exist {
			continue
		}

		if field.Bool() {
			return false
		}
	}

	return true
}
