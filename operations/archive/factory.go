package archive

import "github.com/pufferpanel/pufferpanel/v2"

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Key() string {
	return "archive"
}
func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	source := op.OperationArgs["source"].([]string)
	destination := op.OperationArgs["destination"].(string)

	return Archive{
		Source:      source,
		Destination: destination,
	}, nil
}

var Factory OperationFactory
