package bubblewrap

import (
	"github.com/pufferpanel/pufferpanel/v3"
)

type EnvironmentFactory struct {
	pufferpanel.EnvironmentFactory
}

func (ef EnvironmentFactory) Create() pufferpanel.EnvironmentImpl {
	return &bubblewrap{}
}

func (ef EnvironmentFactory) Key() string {
	return "bubblewrap"
}
