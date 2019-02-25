package view

import (
	"github.com/pufferpanel/pufferpanel/models"
	"gopkg.in/oauth2.v3"
	"strings"
	"time"
)

type OAuthTokenInfoViewModel struct {
	Active   bool                `json:"active"`
	Mapping  map[string][]string `json:"servers,omitempty"`
	Scopes   string              `json:"scope"`
	ClientId string              `json:"client_id"`
}

func FromTokenInfo(info oauth2.TokenInfo, client *models.ClientInfo) *OAuthTokenInfoViewModel {
	model := &OAuthTokenInfoViewModel{}
	model.Active = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn()).After(time.Now())

	if model.Active {
		mapping, scopes := client.MergeServers()

		model.Mapping = mapping
		model.Scopes = strings.Join(scopes, " ")
		model.ClientId = client.ClientID
	}

	return model
}
