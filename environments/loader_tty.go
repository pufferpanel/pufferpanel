//go:build !windows && tty

package environments

import (
	"github.com/pufferpanel/pufferpanel/v3/environments/tty"
)

func init() {
	mapping["standard"] = tty.EnvironmentFactory{}
	mapping["tty"] = tty.EnvironmentFactory{}
}
