package forgedl

import (
	"errors"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/spf13/cast"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	minecraftVersion := cast.ToString(op.OperationArgs["minecraftVersion"])
	version := cast.ToString(op.OperationArgs["version"])
	filename := cast.ToString(op.OperationArgs["target"])

	if version == "" && minecraftVersion == "" {
		return nil, errors.New("missing version and minecraftVersion")
	}

	return ForgeDl{Version: version, Filename: filename, MinecraftVersion: minecraftVersion}, nil
}

func (of OperationFactory) Key() string {
	return "forgedl"
}

var Factory OperationFactory
