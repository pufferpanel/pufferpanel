package extract

import (
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/spf13/cast"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Key() string {
	return "extract"
}
func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	source := cast.ToString(op.OperationArgs["source"])
	destination := cast.ToString(op.OperationArgs["destination"])

	return Extract{
		Source:      source,
		Destination: destination,
	}, nil
}

var Factory OperationFactory
