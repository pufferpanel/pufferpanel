package alterfile

import (
	"bytes"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"os"
	"path/filepath"
	"regexp"
)

type AlterFile struct {
	TargetFile string
	Search     string
	Replace    string
	Regex      bool
}

func (c AlterFile) Run(env pufferpanel.Environment) pufferpanel.OperationResult {
	logging.Info.Printf("Changing data in file: %s", c.TargetFile)
	env.DisplayToConsole(true, "Changing some data in file: %s\n ", c.TargetFile)
	target := filepath.Join(env.GetRootDirectory(), c.TargetFile)
	data, err := os.ReadFile(target)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	var out []byte
	if c.Regex {
		regex, err := regexp.Compile("(?m)" + c.Search)
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}
		out = regex.ReplaceAllLiteral(data, []byte(c.Replace))
	} else {
		out = bytes.ReplaceAll(data, []byte(c.Search), []byte(c.Replace))
	}

	err = os.WriteFile(target, out, 0644)
	return pufferpanel.OperationResult{Error: err}
}
