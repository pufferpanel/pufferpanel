package writefile

import (
	"github.com/pufferpanel/pufferpanel/v2"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Create(op pufferpanel.CreateOperation) pufferpanel.Operation {
	text := op.OperationArgs["text"].(string)
	target := op.OperationArgs["target"].(string)
	return WriteFile{TargetFile: target, Text: text}
}

func (of OperationFactory) Key() string {
	return "writefile"
}

var Factory OperationFactory
