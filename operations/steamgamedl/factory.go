package steamgamedl

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/spf13/cast"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	o := SteamGameDl{
		AppId:     cast.ToString(op.OperationArgs["appId"]),
		Username:  cast.ToString(op.OperationArgs["username"]),
		Password:  cast.ToString(op.OperationArgs["password"]),
		ExtraArgs: cast.ToStringSlice(op.OperationArgs["extraArgs"]),
	}
	return o, nil
}

func (of OperationFactory) Key() string {
	return "steamgamedl"
}

var Factory OperationFactory
