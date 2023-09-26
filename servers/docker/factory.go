package docker

import "github.com/pufferpanel/pufferpanel/v3"

type EnvironmentFactory struct {
	pufferpanel.EnvironmentFactory
}

func (ef EnvironmentFactory) Create(id string) pufferpanel.Environment {
	d := &Docker{
		BaseEnvironment: &pufferpanel.BaseEnvironment{Type: ef.Key(), ServerId: id},
		ContainerId:     id,
		ImageName:       "pufferpanel/generic",
		NetworkMode:     "host",
		Ports:           make([]string, 0),
		Binds:           make(map[string]string),
		Labels:          make(map[string]string),
	}

	d.ExecutionFunction = d.dockerExecuteAsync
	d.BaseEnvironment.WaitFunction = d.WaitForMainProcess
	d.BaseEnvironment.Wrapper = d.CreateWrapper()
	return d
}

func (ef EnvironmentFactory) Key() string {
	return "docker"
}
