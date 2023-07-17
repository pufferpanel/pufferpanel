//go:build !windows && host

package servers

import "github.com/pufferpanel/pufferpanel/v3/servers/tty"

func init() {
	mapping["host"] = tty.EnvironmentFactory{}
	mapping["tty"] = tty.EnvironmentFactory{}
	mapping["standard"] = tty.EnvironmentFactory{}
}
