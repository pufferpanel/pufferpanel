package view

import (
	"gopkg.in/oauth2.v3"
	"time"
)

type OAuthTokenInfoViewModel struct {
	Active  bool                `json:"active"`
	Mapping map[string][]string `json:"mapping,omitempty"`
}

func FromTokenInfo(info oauth2.TokenInfo) *OAuthTokenInfoViewModel {
	model := &OAuthTokenInfoViewModel{}
	model.Active = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn()).After(time.Now())

	//TODO: Copy data from client to mapping ONLY IF ACTIVE
	if model.Active {

	}

	return model
}
