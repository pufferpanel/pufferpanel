package mkdir

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"os"
	"path/filepath"
)

type Mkdir struct {
	TargetFile string
}

func (m *Mkdir) Run(env pufferpanel.Environment) error {
	logging.Info.Printf("Making directory: %s\n", m.TargetFile)
	env.DisplayToConsole(true, "Creating directory: %s\n", m.TargetFile)
	target := filepath.Join(env.GetRootDirectory(), m.TargetFile)
	return os.MkdirAll(target, 0755)
}
