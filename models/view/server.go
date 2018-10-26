package view

import (
	"github.com/pufferpanel/pufferpanel/models"
)

type ServerViewModel struct {
	//Id         uint   `json:"id"`
	Identifier string `json:"id"`
	Name       string `json:"name"`
	NodeId     uint   `json:"nodeId"`
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
