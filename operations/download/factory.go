package download

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/spf13/cast"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	files := cast.ToStringSlice(op.OperationArgs["files"])
	return &Download{Files: files}, nil
}

func (of OperationFactory) Key() string {
	return "download"
}

var Factory OperationFactory
