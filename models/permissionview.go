package models

import "github.com/pufferpanel/pufferpanel/v3"

type PermissionView struct {
	Username         string `json:"username,omitempty"`
	Email            string `json:"email,omitempty"`
	ServerIdentifier string `json:"serverIdentifier,omitempty"`

	Scopes []*pufferpanel.Scope `json:"scopes"`
} //@name Permissions

func FromPermission(p *Permissions) *PermissionView {
	model := &PermissionView{
		Username: p.User.Username,
		Email:    p.User.Email,
		Scopes:   p.Scopes,
	}

	return model
}
