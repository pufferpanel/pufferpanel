package view

import (
	"github.com/pufferpanel/pufferpanel/models"
	"gopkg.in/oauth2.v3"
	"time"
)

type OAuthTokenInfoViewModel struct {
	Active  bool                `json:"active"`
	Mapping map[string][]string `json:"mapping,omitempty"`
}

func FromTokenInfo(info oauth2.TokenInfo, client *models.ClientInfo) *OAuthTokenInfoViewModel {
	model := &OAuthTokenInfoViewModel{}
	model.Active = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn()).After(time.Now())

	if model.Active {
		model.Mapping = client.MergeServers()
	}

	return model
}
