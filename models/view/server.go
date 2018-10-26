package view

import (
	"errors"
	"github.com/pufferpanel/pufferpanel/models"
	"gopkg.in/go-playground/validator.v9"
)

type ServerViewModel struct {
	//Id         uint   `json:"id"`
	Identifier string      `json:"id"`
	Name       string      `json:"name"`
	NodeId     uint        `json:"nodeId"`
	Data       interface{} `json:"data,omitempty"`
}

func FromServer(server *models.Server) *ServerViewModel {
	return &ServerViewModel{
		Name:       server.Name,
		Identifier: server.Identifier,
		NodeId:     server.NodeID,
	}
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

	return nil
}
