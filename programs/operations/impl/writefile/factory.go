package writefile

import "github.com/pufferpanel/pufferpanel/v2/programs/operations/ops"

type OperationFactory struct {
	ops.OperationFactory
}

func (of OperationFactory) Create(op ops.CreateOperation) ops.Operation {
	text := op.OperationArgs["text"].(string)
	target := op.OperationArgs["target"].(string)
	return WriteFile{TargetFile: target, Text: text}
}

func (of OperationFactory) Key() string {
	return "writefile"
}

var Factory OperationFactory
