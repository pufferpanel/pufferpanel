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

import "github.com/pufferpanel/pufferpanel/v3"

type PermissionView struct {
	Username         string `json:"username,omitempty"`
	Email            string `json:"email,omitempty"`
	ServerIdentifier string `json:"serverIdentifier,omitempty"`

	Scopes []pufferpanel.Scope `json:"scopes"`
} //@name Permissions

func FromPermission(p *Permissions) *PermissionView {
	model := &PermissionView{
		Username: p.User.Username,
		Email:    p.User.Email,
		Scopes:   p.Scopes,
	}

	return model
}
