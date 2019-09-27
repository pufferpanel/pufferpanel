package models

import (
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/apufferi/v3/scope"
)

type Permissions struct {
	gorm.Model

	//owners of this permission set
	UserId *uint `json:"-"`
	User   User  `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

	ClientId *uint  `json:"-"`
	Client   Client `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

	//if this set is for a server, what server
	ServerIdentifier *string `json:"-"`
	Server           Server  `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

	//and here are all the perms we support
	Admin bool `gorm:"NOT NULL;DEFAULT 0"`
	//these only will exist if tied to a server
	EditServerData    bool `gorm:"NOT NULL;DEFAULT 0" json:"-"`
	EditServerUsers   bool `gorm:"NOT NULL;DEFAULT 0" json:"-"`
	InstallServer     bool `gorm:"NOT NULL;DEFAULT 0" json:"-"`
	UpdateServer      bool `gorm:"NOT NULL;DEFAULT 0" json:"-"` //this is unused currently
	ViewServerConsole bool `gorm:"NOT NULL;DEFAULT 0" json:"-"`
	SendServerConsole bool `gorm:"NOT NULL;DEFAULT 0" json:"-"`
	StopServer        bool `gorm:"NOT NULL;DEFAULT 0" json:"-"`
	StartServer       bool `gorm:"NOT NULL;DEFAULT 0" json:"-"`
	ViewServerStats   bool `gorm:"NOT NULL;DEFAULT 0" json:"-"`
	ViewServerFiles   bool `gorm:"NOT NULL;DEFAULT 0" json:"-"`
	SFTPServer        bool `gorm:"NOT NULL;DEFAULT 0" json:"-"`
	PutServerFiles    bool `gorm:"NOT NULL;DEFAULT 0" json:"-"`
}

type MultiplePermissions []*Permissions

func (p *Permissions) BeforeSave() {
	if p.ServerIdentifier != nil && *p.ServerIdentifier == "" {
		p.ServerIdentifier = nil
	}
}

func (p *Permissions) ToScopes() []string {
	scopes := make([]string, 0)

	if p.Admin {
		scopes = append(scopes, scope.ServersAdmin)

		if p.ServerIdentifier == nil {
			scopes = append(scopes, scope.ServersCreate, scope.NodesView, scope.NodesDeploy, scope.NodesEdit, scope.TemplatesView, scope.UsersView, scope.UsersEdit)
		} else {
			scopes = append(scopes, scope.ServersDelete, scope.ServersEditAdmin)
		}
	}

	//these only apply if there is a server involved
	if p.ServerIdentifier != nil {
		if p.EditServerData {
			scopes = append(scopes, scope.ServersEdit)
		}

		if p.EditServerUsers {
			scopes = append(scopes, scope.ServersEditUsers)
		}

		if p.InstallServer {
			scopes = append(scopes, scope.ServersInstall)
		}

		if p.UpdateServer {
			//scopes = append(scopes, scope.ServersUpdate)
		}

		if p.ViewServerConsole {
			scopes = append(scopes, scope.ServersConsole)
		}

		if p.SendServerConsole {
			scopes = append(scopes, scope.ServersConsoleSend)
		}

		if p.StopServer {
			scopes = append(scopes, scope.ServersStop)
		}

		if p.StartServer {
			scopes = append(scopes, scope.ServersStart)
		}

		if p.ViewServerStats {
			scopes = append(scopes, scope.ServersStat)
		}

		if p.ViewServerFiles {
			scopes = append(scopes, scope.ServersFilesGet)
		}

		if p.PutServerFiles {
			scopes = append(scopes, scope.ServersFilesPut)
		}

		if p.SFTPServer {
			scopes = append(scopes, scope.ServersFiles)
		}
	}

	return scopes
}
