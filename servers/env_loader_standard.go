//go:build windows && host

package servers

import "github.com/pufferpanel/pufferpanel/v3/servers/standard"

func init() {
	mapping["host"] = standard.EnvironmentFactory{}
	mapping["standard"] = standard.EnvironmentFactory{}
}
