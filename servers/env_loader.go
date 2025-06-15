package servers

import (
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/servers/docker"
	"github.com/pufferpanel/pufferpanel/v3/servers/tty"
	"github.com/pufferpanel/pufferpanel/v3/utils"
	"path/filepath"
	"sync"
)

var envMapping = make(map[string]pufferpanel.EnvironmentFactory)

func init() {
	envMapping["host"] = tty.EnvironmentFactory{}
	envMapping["tty"] = tty.EnvironmentFactory{}
	envMapping["standard"] = tty.EnvironmentFactory{}
	envMapping["docker"] = docker.EnvironmentFactory{}
}

func CreateEnvironment(environmentType, folder string, backupFolder string, server pufferpanel.Server) (*pufferpanel.Environment, error) {
	factory := envMapping[environmentType]

	if factory == nil {
		return nil, fmt.Errorf("undefined environment: %s", environmentType)
	}

	item := &pufferpanel.Environment{
		Type:            factory.Key(),
		ServerId:        server.Identifier,
		ConsoleTracker:  pufferpanel.CreateTracker(),
		StatusTracker:   pufferpanel.CreateTracker(),
		StatsTracker:    pufferpanel.CreateTracker(),
		ConsoleBuffer:   pufferpanel.CreateCache(),
		BackupDirectory: filepath.Join(backupFolder, server.Identifier),
		Wait:            &sync.WaitGroup{},
		Server:          server,
	}
	item.Implementation = factory.Create()
	err := utils.UnmarshalTo(server.Environment.Metadata, item)
	if err != nil {
		return nil, err
	}

	err = utils.UnmarshalTo(server.Environment.Metadata, item.Implementation)
	if err != nil {
		return nil, err
	}

	if item.RootDirectory == "" {
		item.RootDirectory = filepath.Join(folder, server.Identifier)
	}

	item.CreateWrapper()

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
