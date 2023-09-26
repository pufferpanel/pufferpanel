package dockerpull

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/spf13/cast"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	imageName := cast.ToString(op.OperationArgs["imageName"])
	return &DockerPull{ImageName: imageName}, nil
}

func (of OperationFactory) Key() string {
	return "dockerpull"
}

var Factory OperationFactory
