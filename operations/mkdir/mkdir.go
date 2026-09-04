package mkdir

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
)

type Mkdir struct {
	TargetFile string
}

func (m *Mkdir) Run(args pufferpanel.RunOperatorArgs) pufferpanel.OperationResult {
	env := args.Environment
	logging.Info.Printf("Making directory: %s\n", m.TargetFile)
	env.DisplayToConsole(true, "Creating directory: %s\n", m.TargetFile)
	err := args.Server.GetFileServer().MkdirAll(m.TargetFile, 0755)
	return pufferpanel.OperationResult{Error: err}
}
