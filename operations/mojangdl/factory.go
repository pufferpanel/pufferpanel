package mojangdl

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/spf13/cast"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	version := cast.ToString(op.OperationArgs["version"])
	target := cast.ToString(op.OperationArgs["target"])

	return MojangDl{Version: version, Target: target}, nil
}

func (of OperationFactory) Key() string {
	return "mojangdl"
}

var Factory OperationFactory
