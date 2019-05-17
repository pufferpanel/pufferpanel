package view

import (
	"github.com/pufferpanel/pufferpanel/errors"
	"github.com/pufferpanel/pufferpanel/models"
	"gopkg.in/go-playground/validator.v9"
	"net/url"
)

type NodeViewModel struct {
	Id          uint   `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	PublicHost  string `json:"publicHost,omitempty"`
	PrivateHost string `json:"privateHost,omitempty"`
	PublicPort  uint   `json:"publicPort,omitempty"`
	PrivatePort uint   `json:"privatePort,omitempty"`
	SFTPPort    uint   `json:"sftpPort,omitempty"`
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

func (n *NodeViewModel) Valid(allowEmpty bool) error {
	validate := validator.New()

	if !allowEmpty && validate.Var(n.Name, "required") != nil {
		return errors.ErrFieldRequired("name")
	}

	if validate.Var(n.Name, "optional|printascii") != nil {
		return errors.ErrFieldMustBePrintable("name")
	}

	testName := url.QueryEscape(n.Name)
	if testName != n.Name {
		return errors.ErrFieldHasURICharacters("name")
	}

	if !allowEmpty && validate.Var(n.PublicHost, "required") != nil {
		return errors.ErrFieldMustBePrintable("publicHost")
	}

	if validate.Var(n.PublicHost, "optional|ip|fqdn") != nil {
		return errors.ErrFieldIsInvalidHost("publicHost")
	}

	if !allowEmpty && validate.Var(n.PrivateHost, "required") != nil {
		return errors.ErrFieldMustBePrintable("privateHost")
	}

	if validate.Var(n.PrivateHost, "optional|ip_addr|fqdn") != nil {
		return errors.ErrFieldIsInvalidHost("privateHost")
	}

	if allowEmpty {
		if validate.Var(n.PublicPort, "min=0,max=65535") != nil {
			return errors.ErrFieldTooLarge("publicPort", 65535)
		}

		if validate.Var(n.PrivatePort, "min=0,max=65535") != nil {
			return errors.ErrFieldTooLarge("privatePort", 65535)
		}

		if validate.Var(n.SFTPPort, "min=0,max=65535") != nil {
			return errors.ErrFieldTooLarge("sftpPort", 65535)
		}
	} else {
		if validate.Var(n.PublicPort, "min=1,max=65535") != nil {
			return errors.ErrFieldNotBetween("publicPort", 1, 65535)
		}

		if validate.Var(n.PrivatePort, "min=1,max=65535") != nil {
			return errors.ErrFieldNotBetween("privatePort", 1, 65535)
		}

		if validate.Var(n.SFTPPort, "min=1,max=65535") != nil {
			return errors.ErrFieldNotBetween("sftpPort", 1, 65535)
		}
	}

	if n.SFTPPort != 0 && n.SFTPPort == n.PublicPort {
		return errors.ErrFieldEqual("sftpPort", "publicPort")
	}

	if n.SFTPPort != 0 && n.SFTPPort == n.PrivatePort {
		return errors.ErrFieldEqual("sftpPort", "privatePort")
	}

	return nil
}
