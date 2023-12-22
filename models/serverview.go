package models

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"gopkg.in/go-playground/validator.v9"
)

type ServerView struct {
	Identifier   string           `json:"id,omitempty"`
	Name         string           `json:"name,omitempty"`
	NodeId       uint             `json:"nodeId,omitempty"`
	Node         *NodeView        `json:"node,omitempty"`
	Data         interface{}      `json:"data,omitempty"`
	Users        []ServerUserView `json:"users,omitempty"`
	IP           string           `json:"ip,omitempty"`
	Port         uint16           `json:"port,omitempty"`
	Type         string           `json:"type"`
	Icon         string           `json:"icon,omitempty"`
	CanGetStatus bool             `json:"canGetStatus,omitempty"`
} //@name ServerInfo

type ServerUserView struct {
	Username string   `json:"username"`
	Scopes   []string `json:"scopes"`
} //@name ServerUser

func FromServer(server *Server) *ServerView {
	model := &ServerView{
		Name:       server.Name,
		Identifier: server.Identifier,
		NodeId:     server.NodeID,
		IP:         server.IP,
		Port:       server.Port,
		Type:       server.Type,
		Icon:       server.Icon,
		Node:       FromNode(&server.Node),
	}

	return model
}

func FromServers(servers []*Server) []*ServerView {
	result := make([]*ServerView, len(servers))

	for k, v := range servers {
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

	if validate.Var(s.Name, "omitempty,printascii") != nil {
		return pufferpanel.ErrFieldMustBePrintable("name")
	}

	if !allowEmpty && validate.Var(s.NodeId, "required,min:1") != nil {
		return pufferpanel.ErrFieldTooSmall("node", 1)
	}

	if validate.Var(s.IP, "omitempty,ip|fqdn") != nil {
		return pufferpanel.ErrFieldIsInvalidIP("ip")
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
