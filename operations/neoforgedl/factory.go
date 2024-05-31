package neoforgedl

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
	outputVariable := cast.ToString(op.OperationArgs["outputVariable"])

	if version == "" && minecraftVersion == "" {
		return nil, errors.New("missing version and minecraftVersion")
	}

	if outputVariable == "" {
		outputVariable = "opForgeVersion"
	}

	return neoForgeDl{Version: version, Filename: filename, MinecraftVersion: minecraftVersion, OutputVariable: outputVariable}, nil
}

func (of OperationFactory) Key() string {
	return "neoforgedl"
}

var Factory OperationFactory
