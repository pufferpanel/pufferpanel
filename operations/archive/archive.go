package archive

import (
	"github.com/mholt/archiver"
	"github.com/pufferpanel/pufferpanel/v2"
)

type Archive struct {
	Source      []string
	Destination string
}

func (op Archive) Run(pufferpanel.Environment) error {
	return archiver.Archive(op.Source, op.Destination)
}
