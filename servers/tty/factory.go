//go:build !windows

package tty

import "github.com/pufferpanel/pufferpanel/v3"

type EnvironmentFactory struct {
	pufferpanel.EnvironmentFactory
}

func (ef EnvironmentFactory) Create(id string) pufferpanel.Environment {
	t := &tty{
		BaseEnvironment: &pufferpanel.BaseEnvironment{Type: ef.Key(), ServerId: id},
	}
	t.BaseEnvironment.ExecutionFunction = t.ttyExecuteAsync
	t.BaseEnvironment.WaitFunction = t.WaitForMainProcess
	t.BaseEnvironment.Wrapper = t.CreateWrapper()
	return t
}

func (ef EnvironmentFactory) Key() string {
	return "host"
}
