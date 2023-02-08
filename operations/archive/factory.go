package archive

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/spf13/cast"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Key() string {
	return "archive"
}
func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	source := cast.ToStringSlice(op.OperationArgs["source"])
	destination := cast.ToString(op.OperationArgs["destination"])

	return Archive{
		Source:      source,
		Destination: destination,
	}, nil
}

var Factory OperationFactory
