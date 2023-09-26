package mkdir

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/spf13/cast"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	target := cast.ToString(op.OperationArgs["target"])
	return &Mkdir{TargetFile: target}, nil
}

func (of OperationFactory) Key() string {
	return "mkdir"
}

var Factory OperationFactory
