/*
 Copyright 2019 Padduck, LLC
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
	"gopkg.in/oauth2.v3"
	"strings"
	"time"
)

type TokenInfoView struct {
	Active   bool                `json:"active"`
	Mapping  map[string][]string `json:"servers,omitempty"`
	Scopes   string              `json:"scope"`
	ClientId string              `json:"client_id"`
}

func FromTokenInfo(info oauth2.TokenInfo, client *ClientInfo) *TokenInfoView {
	model := &TokenInfoView{}
	model.Active = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn()).After(time.Now())

	if model.Active {
		mapping, scopes := client.MergeServers()

		model.Mapping = mapping
		model.Scopes = strings.Join(scopes, " ")
		model.ClientId = client.ClientID
	}

	return model
}
