package paperdl

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
	build := cast.ToString(op.OperationArgs["build"])
	filename := cast.ToString(op.OperationArgs["target"])
	project := cast.ToString(op.OperationArgs["project"])

	if minecraftVersion == "" {
		return nil, errors.New("missing minecraftVersion")
	}

	if build == "" {
		return nil, errors.New("missing build")
	}

	if filename == "" {
		return nil, errors.New("missing target")
	}

	if project == "" {
		project = "paper"
	}

	return PaperDl{MinecraftVersion: minecraftVersion, Build: build, Filename: filename, Project: project}, nil
}

func (of OperationFactory) Key() string {
	return "paperdl"
}

var Factory OperationFactory
