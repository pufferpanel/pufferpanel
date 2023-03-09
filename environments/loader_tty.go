//go:build !windows && host

package environments

import (
	"github.com/pufferpanel/pufferpanel/v3/environments/tty"
)

func init() {
	mapping["host"] = tty.EnvironmentFactory{}
	mapping["tty"] = tty.EnvironmentFactory{}
	mapping["standard"] = tty.EnvironmentFactory{}
}
