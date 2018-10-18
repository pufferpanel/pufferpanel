package view

import "github.com/pufferpanel/pufferpanel/models"

type NodeViewModel struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	PublicHost  string `json:"publicHost"`
	PrivateHost string `json:"privateHost"`
	PublicPort  int    `json:"publicPort"`
	PrivatePort int    `json:"privatePort"`
	SFTPPort    int    `json:"sftpPort"`
}

func FromNode (n models.Node) NodeViewModel {
	return NodeViewModel{
		Id: n.Id,
		Name: n.Name,
		PublicHost: n.PublicHost,
		PrivateHost: n.PrivateHost,
		PublicPort: n.PublicPort,
		PrivatePort: n.PrivatePort,
		SFTPPort: n.SFTPPort,
	}
}

func (n NodeViewModel) CopyToModel(newModel models.Node) {
	if n.Name != "" {
		newModel.Name = n.Name
	}

	if n.PublicHost != "" {
		newModel.PublicHost = n.PublicHost
	}

	if n.PrivateHost != "" {
		newModel.PrivateHost = n.PrivateHost
	}

	if n.PublicPort > 0 {
		newModel.PublicPort = n.PublicPort
	}

	if n.PrivatePort > 0 {
		newModel.PrivatePort = n.PrivatePort
	}

	if n.SFTPPort > 0 {
		newModel.SFTPPort = n.SFTPPort
	}
}