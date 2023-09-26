package dockerpull

import (
	"context"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/servers/docker"
)

type DockerPull struct {
	ImageName string
}

func (d DockerPull) Run(env pufferpanel.Environment) error {
	dockerEnv, ok := env.(*docker.Docker)

	if !ok {
		return pufferpanel.ErrEnvironmentNotSupported
	}

	return dockerEnv.PullImage(context.Background(), d.ImageName, true)
}
