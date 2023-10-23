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
		Network:         "host",
		Ports:           make([]string, 0),
		Binds:           make(map[string]string),
		Labels:          make(map[string]string),
	}

	d.ExecutionFunction = d.dockerExecuteAsync
	d.BaseEnvironment.Wrapper = d.CreateWrapper()
	d.BaseEnvironment.IsRunningFunc = d.IsRunning
	return d
}

func (ef EnvironmentFactory) Key() string {
	return "docker"
}
