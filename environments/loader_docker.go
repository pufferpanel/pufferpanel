//go:build docker

package environments

import "github.com/pufferpanel/pufferpanel/v3/environments/docker"

func init() {
	mapping["docker"] = docker.EnvironmentFactory{}
}
