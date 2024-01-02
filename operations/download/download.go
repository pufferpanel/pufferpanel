package download

import (
	"github.com/cavaliergopher/grab/v3"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
)

type Download struct {
	Files []string
}

func (d Download) Run(env pufferpanel.Environment) pufferpanel.OperationResult {
	for _, file := range d.Files {
		logging.Info.Printf("Download file from %s to %s", file, env.GetRootDirectory())
		env.DisplayToConsole(true, "Downloading file %s\n", file)
		_, err := grab.Get(env.GetRootDirectory(), file)
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}
	}
	return pufferpanel.OperationResult{Error: nil}
}
