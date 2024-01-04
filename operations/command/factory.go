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

	var stdIn pufferpanel.StdinConsoleConfiguration
	if field, exists := op.OperationArgs["stdin"]; exists {
		err := pufferpanel.UnmarshalTo(field, stdIn)
		if err != nil {
			return nil, err
		}
	}

	return Command{Commands: cmds, Env: op.EnvironmentVariables, StdIn: stdIn, Variables: op.DataMap}, nil
}

func (of OperationFactory) Key() string {
	return "command"
}

var Factory OperationFactory
