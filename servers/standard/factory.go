package standard

import "github.com/pufferpanel/pufferpanel/v3"

type EnvironmentFactory struct {
	pufferpanel.EnvironmentFactory
}

func (ef EnvironmentFactory) Create(id string) pufferpanel.Environment {
	s := &standard{
		BaseEnvironment: &pufferpanel.BaseEnvironment{Type: ef.Key(), ServerId: id},
	}

	s.BaseEnvironment.Wrapper = s.CreateWrapper()
	s.BaseEnvironment.ExecutionFunction = s.standardExecuteAsync
	s.BaseEnvironment.IsRunningFunc = s.isRunning
	s.BaseEnvironment.KillFunc = s.kill
	return s
}

func (ef EnvironmentFactory) Key() string {
	return "host"
}
