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
	"github.com/pufferpanel/pufferpanel"
	"gopkg.in/go-playground/validator.v9"
)

type ServerView struct {
	Identifier string           `json:"id,omitempty"`
	Name       string           `json:"name,omitempty"`
	NodeId     uint             `json:"nodeId,omitempty"`
	Node       *NodeView        `json:"node,omitempty"`
	Data       interface{}      `json:"data,omitempty"`
	Users      []ServerUserView `json:"users,omitempty"`
	IP         string           `json:"ip,omitempty"`
	Port       uint             `json:"port,omitempty"`
	Type       string           `json:"type"`
}

type ServerUserView struct {
	Username string   `json:"username"`
	Scopes   []string `json:"scopes"`
}

func FromServer(server *Server) *ServerView {
	model := &ServerView{
		Name:       server.Name,
		Identifier: server.Identifier,
		NodeId:     server.NodeID,
		IP:         server.IP,
		Port:       server.Port,
		Type:       server.Type,
	}

	if server.Node.ID != 0 {
		model.Node = FromNode(&server.Node)
	}

	return model
}

func FromServers(servers *Servers) []*ServerView {
	result := make([]*ServerView, len(*servers))

	for k, v := range *servers {
		result[k] = FromServer(v)
	}

	return result
}

func (s *ServerView) Valid(allowEmpty bool) error {
	validate := validator.New()

	if !allowEmpty && validate.Var(s.Name, "required") != nil {
		return pufferpanel.ErrFieldRequired("name")
	}

	if !allowEmpty && validate.Var(s.Type, "required") != nil {
		return pufferpanel.ErrFieldRequired("type")
	}

	if validate.Var(s.Name, "optional|printascii") != nil {
		return pufferpanel.ErrFieldMustBePrintable("name")
	}

	if !allowEmpty && validate.Var(s.NodeId, "required,min:1") != nil {
		return pufferpanel.ErrFieldTooSmall("node", 1)
	}

	if validate.Var(s.IP, "optional|ip_addr") != nil {
		return pufferpanel.ErrFieldIsInvalidIP("ip")
	}

	if validate.Var(s.Port, "optional|min:0,max:65535") != nil {
		return pufferpanel.ErrFieldNotBetween("port", 1, 65535)
	}

	return nil
}

func RemoveServerPrivateInfoFromAll(servers []*ServerView) []*ServerView {
	for k, v := range servers {
		servers[k] = RemoveServerPrivateInfo(v)
	}
	return servers
}

func RemoveServerPrivateInfo(server *ServerView) *ServerView {
	//SCRUB DATA FROM REGULAR USERS
	if server.Node != nil {
		server.Node.Id = 0
		server.NodeId = 0
		server.Node.PrivateHost = ""
		server.Node.PrivatePort = 0
	}

	return server
}
