package writefile

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"os"
)

type WriteFile struct {
	TargetFile string
	Text       string
}

func (c WriteFile) Run(args pufferpanel.RunOperatorArgs) pufferpanel.OperationResult {
	env := args.Environment
	fs := args.Server.GetFileServer()

	logging.Info.Printf("Writing data to file: %s", c.TargetFile)
	env.DisplayToConsole(true, "Writing some data to file: %s\n", c.TargetFile)

	file, err := fs.OpenFile(c.TargetFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}
	defer pufferpanel.Close(file)

	_, err = file.Write([]byte(c.Text))
	return pufferpanel.OperationResult{Error: err}
}
