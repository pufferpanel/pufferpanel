package command

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/spf13/cast"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	cmds := cast.ToStringSlice(op.OperationArgs["commands"])
	return Command{Commands: cmds, Env: op.EnvironmentVariables}, nil
}

func (of OperationFactory) Key() string {
	return "command"
}

var Factory OperationFactory
