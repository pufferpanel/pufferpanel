package archive

import (
	"github.com/mholt/archiver/v3"
	"github.com/pufferpanel/pufferpanel/v3"
)

type Archive struct {
	Source      []string
	Destination string
}

func (op Archive) Run(pufferpanel.Environment) pufferpanel.OperationResult {
	err := archiver.Archive(op.Source, op.Destination)
	return pufferpanel.OperationResult{Error: err}
}
