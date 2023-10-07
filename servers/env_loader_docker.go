//go:build !nodocker

package servers

import "github.com/pufferpanel/pufferpanel/v3/servers/docker"

func init() {
	envMapping["docker"] = docker.EnvironmentFactory{}
}
