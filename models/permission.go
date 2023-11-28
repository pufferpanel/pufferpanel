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

	RawScopes string               `gorm:"column:scopes;NOT NULL" json:"-" validate:"required"`
	Scopes    []*pufferpanel.Scope `gorm:"-" json:"-"`
}

func (p *Permissions) BeforeSave(*gorm.DB) error {
	if p.ServerIdentifier != nil && *p.ServerIdentifier == "" {
		p.ServerIdentifier = nil
	}

	if p.ServerIdentifier != nil {
		//ensure they have the view, because we're saving them back in
		p.Scopes = pufferpanel.AddScope(p.Scopes, pufferpanel.ScopeServerView)
	}

	tmp := make([]string, len(p.Scopes))
	for k, v := range p.Scopes {
		tmp[k] = v.String()
	}
	p.RawScopes = strings.Join(tmp, ",")
	return nil
}

func (p *Permissions) AfterFind(*gorm.DB) error {
	p.Scopes = make([]*pufferpanel.Scope, 0)
	if p.RawScopes != "" {
		for _, v := range strings.Split(p.RawScopes, ",") {
			//we can just simply blindly assign it, because the checks we do are just making these strings anyways...
			p.Scopes = append(p.Scopes, pufferpanel.GetScope(v))
		}
	}

	return nil
}

func (p *Permissions) ShouldDelete() bool {
	if p.ServerIdentifier == nil {
		return false
	}
	return len(p.Scopes) == 0
}
