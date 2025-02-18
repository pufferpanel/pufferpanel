package stdin

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/spf13/cast"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	message := cast.ToString(op.OperationArgs["command"])
	return &Stdin{Command: message}, nil
}

func (of OperationFactory) Key() string {
	return "stdin"
}

var Factory OperationFactory
