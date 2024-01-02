package dockerpull

import (
	"context"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/servers/docker"
)

type DockerPull struct {
	ImageName string
}

func (d DockerPull) Run(env pufferpanel.Environment) pufferpanel.OperationResult {
	dockerEnv, ok := env.(*docker.Docker)

	if !ok {
		return pufferpanel.OperationResult{Error: pufferpanel.ErrEnvironmentNotSupported}
	}

	err := dockerEnv.PullImage(context.Background(), d.ImageName, true)
	return pufferpanel.OperationResult{Error: err}
}
