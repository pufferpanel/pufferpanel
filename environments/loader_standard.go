//go:build windows && host

package environments

import "github.com/pufferpanel/pufferpanel/v3/environments/standard"

func init() {
	mapping["host"] = standard.EnvironmentFactory{}
}
