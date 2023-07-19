//go:build docker

package servers

import "github.com/pufferpanel/pufferpanel/v3/servers/docker"

func init() {
	envMapping["docker"] = docker.EnvironmentFactory{}
}
