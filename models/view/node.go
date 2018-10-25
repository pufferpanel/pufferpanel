package view

import (
	"errors"
	"github.com/pufferpanel/pufferpanel/models"
	"gopkg.in/go-playground/validator.v9"
	"net/url"
)

type NodeViewModel struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	PublicHost  string `json:"publicHost"`
	PrivateHost string `json:"privateHost"`
	PublicPort  uint   `json:"publicPort"`
	PrivatePort uint   `json:"privatePort"`
	SFTPPort    uint   `json:"sftpPort"`
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

func (n *NodeViewModel) Valid() error {
	validate := validator.New()

	if validate.Var(n.Name, "required") != nil {
		return errors.New("name is required")
	}

	if validate.Var(n.Name, "printascii") != nil {
		return errors.New("name must be printable ascii characters")
	}

	testName := url.QueryEscape(n.Name)
	if testName != n.Name {
		return errors.New("name must not contain characters which cannot be used in URIs")
	}

	if validate.Var(n.PublicHost, "required") != nil {
		return errors.New("publicHost is required")
	}

	if validate.Var(n.PublicHost, "ip|fqdn") != nil {
		return errors.New("publicHost must be a valid IP or FQDN")
	}

	if validate.Var(n.PrivateHost, "required") != nil {
		return errors.New("privateHost is required")
	}

	if validate.Var(n.PrivateHost, "ip_addr|fqdn") != nil {
		return errors.New("privateHost must be a valid IP or FQDN")
	}

	if validate.Var(n.PublicPort, "min=1,max=65535") != nil {
		return errors.New("publicPort must be between 1 and 65535")
	}

	if validate.Var(n.PrivatePort, "min=1,max=65535") != nil {
		return errors.New("privatePort must be between 1 and 65535")
	}

	if validate.Var(n.SFTPPort, "min=1,max=65535") != nil {
		return errors.New("sftpPort must be between 1 and 65535")
	}

	if n.SFTPPort == n.PublicPort {
		return errors.New("sftpPort cannot be the same as the public port")
	}

	if n.SFTPPort == n.PrivatePort {
		return errors.New("sftpPort cannot be the same as the private port")
	}

	return nil
}
