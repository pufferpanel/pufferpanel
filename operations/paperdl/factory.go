package paperdl

import (
	"errors"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/spf13/cast"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	logging.Debug.Printf("create entered")
	minecraftVersion := cast.ToString(op.OperationArgs["minecraftVersion"])
	build := cast.ToString(op.OperationArgs["build"])
	filename := cast.ToString(op.OperationArgs["target"])

	if minecraftVersion == "" {
		return nil, errors.New("missing minecraftVersion")
	}

	if build == "" {
		return nil, errors.New("missing build")
	}

	logging.Debug.Printf("create done")
	return PaperDl{MinecraftVersion: minecraftVersion, Build: build, Filename: filename}, nil
}

func (of OperationFactory) Key() string {
	return "paperdl"
}

var Factory OperationFactory
