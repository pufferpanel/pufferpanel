//go:build !windows && host

package servers

import "github.com/pufferpanel/pufferpanel/v3/servers/tty"

func init() {
	envMapping["host"] = tty.EnvironmentFactory{}
	envMapping["tty"] = tty.EnvironmentFactory{}
	envMapping["standard"] = tty.EnvironmentFactory{}
}
