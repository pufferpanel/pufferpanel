package view

import "github.com/pufferpanel/pufferpanel/models"

type NodeViewModel struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	PublicHost  string `json:"publicHost"`
	PrivateHost string `json:"privateHost"`
	PublicPort  uint    `json:"publicPort"`
	PrivatePort uint    `json:"privatePort"`
	SFTPPort    uint    `json:"sftpPort"`
}

func FromNode(n *models.Node) *NodeViewModel {
	return &NodeViewModel{
		Id:          n.ID,
		Name:        n.Name,
		PublicHost:  n.PublicHost,
		PrivateHost: n.PrivateHost,
		PublicPort:  n.PublicPort,
		PrivatePort: n.PrivatePort,
		SFTPPort:    n.SFTPPort,
	}
}

func FromNodes(n *models.Nodes) []*NodeViewModel {
	result := make([]*NodeViewModel, len(*n))

	for k, v := range *n {
		result[k] = FromNode(v)
	}

	return result
}

func (n *NodeViewModel) CopyToModel(newModel *models.Node) {
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
