package extract

import (
	"github.com/mholt/archiver/v3"
	"github.com/pufferpanel/pufferpanel/v3"
)

type Extract struct {
	Source      string
	Destination string
}

func (op Extract) Run(pufferpanel.Environment) error {
	return archiver.Unarchive(op.Source, op.Destination)
}
