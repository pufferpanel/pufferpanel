package servers

import (
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"path/filepath"
	"sync"
)

var envMapping = make(map[string]pufferpanel.EnvironmentFactory)

func CreateEnvironment(environmentType, folder, id string, environmentSection pufferpanel.MetadataType) (pufferpanel.Environment, error) {
	factory := envMapping[environmentType]

	if factory == nil {
		switch environmentType {
		case "standard":
			factory = envMapping["host"]
		case "tty":
			factory = envMapping["host"]
		}
	}

	if factory == nil {
		return nil, fmt.Errorf("undefined environment: %s", environmentType)
	}

	item := factory.Create(id)
	err := pufferpanel.UnmarshalTo(environmentSection.Metadata, item)
	if err != nil {
		return nil, err
	}

	serverRoot := filepath.Join(folder, id)
	envCache := pufferpanel.CreateCache()

	e := item.GetBase()
	if e.RootDirectory == "" {
		e.RootDirectory = serverRoot
	}
	e.ConsoleTracker = pufferpanel.CreateTracker()
	e.StatusTracker = pufferpanel.CreateTracker()
	e.StatsTracker = pufferpanel.CreateTracker()

	e.ConsoleBuffer = envCache
	e.Wait = &sync.WaitGroup{}
	e.Wrapper = e.CreateWrapper()

	return item, nil
}

func GetSupportedEnvironments() []string {
	deduper := make(map[string]bool)

	for k := range envMapping {
		deduper[k] = true
	}

	result := make([]string, len(deduper))
	i := 0
	for k := range deduper {
		result[i] = k
		i++
	}

	return result
}
