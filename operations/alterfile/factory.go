package alterfile

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/spf13/cast"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	file := cast.ToString(op.OperationArgs["file"])
	search := cast.ToString(op.OperationArgs["search"])
	replace := cast.ToString(op.OperationArgs["replace"])
	regex := cast.ToBool(op.OperationArgs["regex"])
	return AlterFile{TargetFile: file, Search: search, Replace: replace, Regex: regex}, nil
}

func (of OperationFactory) Key() string {
	return "alterfile"
}

var Factory OperationFactory
