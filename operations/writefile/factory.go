package writefile

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/spf13/cast"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	text := cast.ToString(op.OperationArgs["text"])
	target := cast.ToString(op.OperationArgs["target"])
	return WriteFile{TargetFile: target, Text: text}, nil
}

func (of OperationFactory) Key() string {
	return "writefile"
}

var Factory OperationFactory
