package dockerpull

import (
	"context"

	"github.com/pufferpanel/pufferpanel/v3"
)

type DockerPull struct {
	ImageName string
}

func (d DockerPull) Run(args pufferpanel.RunOperatorArgs) pufferpanel.OperationResult {
	env := args.Environment
	err := pufferpanel.PullDockerImage(env, context.Background(), d.ImageName, true)
	return pufferpanel.OperationResult{Error: err}
}
