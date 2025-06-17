package docker

import (
	"github.com/pufferpanel/pufferpanel/v3"
)

type EnvironmentFactory struct {
	pufferpanel.EnvironmentFactory
}

func (ef EnvironmentFactory) Create() pufferpanel.EnvironmentImpl {
	return &Docker{
		ImageName: "pufferpanel/generic",
		Network:   "host",
		Ports:     make([]string, 0),
		Binds:     make(map[string]string),
		Labels:    make(map[string]string),
	}
}

func (ef EnvironmentFactory) Key() string {
	return "docker"
}
