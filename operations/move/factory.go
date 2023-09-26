package move

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/spf13/cast"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	source := cast.ToString(op.OperationArgs["source"])
	target := cast.ToString(op.OperationArgs["target"])
	return Move{SourceFile: source, TargetFile: target}, nil
}

func (of OperationFactory) Key() string {
	return "move"
}

var Factory OperationFactory
