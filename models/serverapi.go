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
	"github.com/pufferpanel/pufferpanel/v2/response"
)

type ServerCreation struct {
	pufferpanel.Server

	NodeId uint     `json:"node"`
	Users  []string `json:"users"`
	Name   string   `json:"name"`
}

type GetServerResponse struct {
	Server *ServerView     `json:"server"`
	Perms  *PermissionView `json:"permissions"`
}

type CreateServerResponse struct {
	Id string `json:"id"`
}

type ServerSearchResponse struct {
	Servers []*ServerView `json:"servers"`
	*response.Metadata
}
