package standard

import "github.com/pufferpanel/pufferpanel/v3"

type EnvironmentFactory struct {
	pufferpanel.EnvironmentFactory
}

func (ef EnvironmentFactory) Create(id string) pufferpanel.Environment {
	s := &standard{
		BaseEnvironment: &pufferpanel.BaseEnvironment{Type: ef.Key(), ServerId: id},
	}
	s.BaseEnvironment.ExecutionFunction = s.standardExecuteAsync
	s.BaseEnvironment.Wrapper = s.CreateWrapper()
	s.BaseEnvironment.IsRunningFunc = s.IsRunning
	return s
}

func (ef EnvironmentFactory) Key() string {
	return "host"
}
