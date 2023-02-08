//go:build windows && standard

package environments

import "github.com/pufferpanel/pufferpanel/v3/environments/standard"

func init() {
	mapping["standard"] = standard.EnvironmentFactory{}
}
