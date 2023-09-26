package spongedl

import (
	"encoding/json"
	"github.com/pufferpanel/pufferpanel/v3"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	jsonData, err := json.Marshal(op.OperationArgs)
	if err != nil {
		return nil, err
	}

	var spongedl SpongeDl
	err = json.Unmarshal(jsonData, &spongedl)
	return spongedl, err
}

func (of OperationFactory) Key() string {
	return "spongedl"
}

var Factory OperationFactory
