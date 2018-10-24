package view

import (
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/satori/go.uuid"
)

type ServerViewModel struct {
	Id     uint      `json:"id"`
	Name   string    `json:"name"`
	UUID   uuid.UUID `json:"uuid"`
	NodeId uint      `json:"nodeId"`
}

func FromServer(server *models.Server) *ServerViewModel {
	return &ServerViewModel{
		Id:     server.ID,
		Name:   server.Name,
		UUID:   server.UUID,
		NodeId: server.NodeID,
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
