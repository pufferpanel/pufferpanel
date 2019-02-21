package view

import (
	"errors"
	"github.com/pufferpanel/pufferpanel/models"
	"gopkg.in/go-playground/validator.v9"
)

type ServerViewModel struct {
	//Id         uint   `json:"id"`
	Identifier string                `json:"id,omitempty"`
	Name       string                `json:"name,omitempty"`
	NodeId     uint                  `json:"nodeId,omitempty"`
	Node       *NodeViewModel        `json:"node,omitempty"`
	Data       interface{}           `json:"data,omitempty"`
	Users      []ServerViewModelUser `json:"users,omitempty"`
	IP         string                `json:"ip,omitempty"`
	Port       uint                  `json:"port,omitempty"`
}

type ServerViewModelUser struct {
	Username string   `json:"username"`
	Scopes   []string `json:"scopes"`
}

func FromServer(server *models.Server) *ServerViewModel {
	model := &ServerViewModel{
		Name:       server.Name,
		Identifier: server.Identifier,
		NodeId:     server.NodeID,
		IP:         server.IP,
		Port:       server.Port,
	}

	if server.Node.ID != 0 {
		model.Node = FromNode(&server.Node)
	}

	return model
}

func FromServers(servers *models.Servers) []*ServerViewModel {
	result := make([]*ServerViewModel, len(*servers))

	for k, v := range *servers {
		result[k] = FromServer(v)
	}

	return result
}

func (model *ServerViewModel) CopyToModel(newModel *models.Server) {
	if model.Name != "" {
		newModel.Name = model.Name
	}
}

func (s *ServerViewModel) Valid(allowEmpty bool) error {
	validate := validator.New()

	if !allowEmpty && validate.Var(s.Name, "required") != nil {
		return errors.New("name is required")
	}

	if validate.Var(s.Name, "optional|printascii") != nil {
		return errors.New("name must be printable ascii characters")
	}

	if !allowEmpty && validate.Var(s.NodeId, "required,min:1") != nil {
		return errors.New("node id must be a positive non-zero number")
	}

	if !allowEmpty && validate.Var(s.IP, "optional|min:0,max:65535") != nil {
		return errors.New("port must either not be included or be between 0 and 65535")
	}

	return nil
}

func RemoveServerPrivateInfoFromAll(servers []*ServerViewModel) []*ServerViewModel {
	for k, v := range servers {
		servers[k] = RemoveServerPrivateInfo(v)
	}
	return servers
}

func RemoveServerPrivateInfo(server *ServerViewModel) *ServerViewModel {
	//SCRUB DATA FROM REGULAR USERS
	if server.Node != nil {
		server.Node.Id = 0
		server.NodeId = 0
		server.Node.PrivateHost = ""
		server.Node.PrivatePort = 0
	}

	return server
}
