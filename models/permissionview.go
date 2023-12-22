package models

import "github.com/pufferpanel/pufferpanel/v3"

type PermissionView struct {
	ServerIdentifier string `json:"serverIdentifier,omitempty"`

	Scopes []*pufferpanel.Scope `json:"scopes"`
} //@name Permissions

func FromPermission(p *Permissions) *PermissionView {
	model := &PermissionView{
		Scopes: p.Scopes,
	}

	if model.Scopes == nil {
		model.Scopes = make([]*pufferpanel.Scope, 0)
	}

	return model
}

type UserPermissionsView struct {
	Username string               `json:"username,omitempty"`
	Email    string               `json:"email"`
	Scopes   []*pufferpanel.Scope `json:"scopes"`
}
