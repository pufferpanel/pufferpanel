package writefile

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"os"
	"path/filepath"
)

type WriteFile struct {
	TargetFile string
	Text       string
}

func (c WriteFile) Run(env pufferpanel.Environment) pufferpanel.OperationResult {
	logging.Info.Printf("Writing data to file: %s", c.TargetFile)
	env.DisplayToConsole(true, "Writing some data to file: %s\n", c.TargetFile)
	target := filepath.Join(env.GetRootDirectory(), c.TargetFile)
	err := os.WriteFile(target, []byte(c.Text), 0644)
	return pufferpanel.OperationResult{Error: err}
}
