package view

import (
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/satori/go.uuid"
)

type ServerViewModel struct {
	Id     int       `json:"id"`
	Name   string    `json:"name"`
	UUID   uuid.UUID `json:"uuid"`
	NodeId int       `json:"nodeId"`
}

func FromServer(server *models.Server) *ServerViewModel {
	return &ServerViewModel{
		Id:     server.Id,
		Name:   server.Name,
		UUID:   server.UUID,
		NodeId: server.NodeId,
	}
}

func (model *ServerViewModel) CopyToModel(newModel *models.Server) {
	if model.Name != "" {
		newModel.Name = model.Name
	}
}
