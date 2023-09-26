package javadl

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/spf13/cast"
	"strconv"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Key() string {
	return "javadl"
}
func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	v := op.OperationArgs["version"]

	version, err := cast.ToIntE(v)
	if err != nil {
		return nil, err
	}

	return JavaDl{Version: strconv.Itoa(version)}, nil
}

var Factory OperationFactory
