package sleep

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/spf13/cast"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	duration, err := cast.ToDurationE(op.OperationArgs["duration"])
	if err != nil {
		return nil, err
	}
	return &Sleep{Duration: duration}, nil
}

func (of OperationFactory) Key() string {
	return "sleep"
}

var Factory OperationFactory
