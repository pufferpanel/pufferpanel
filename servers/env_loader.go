package servers

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/utils"
)

var envMapping = make(map[string]func() pufferpanel.EnvironmentImpl)

func init() {
	envMapping["host"] = CreateTTYEnvironment
	envMapping["tty"] = CreateTTYEnvironment
	envMapping["standard"] = CreateTTYEnvironment
	envMapping["docker"] = CreateDockerEnvironment
}

func CreateEnvironment(environmentType, folder string, backupFolder string, server *Server) (*pufferpanel.Environment, error) {
	factory := envMapping[environmentType]

	if factory == nil {
		return nil, fmt.Errorf("undefined environment: %s", environmentType)
	}

	item := &pufferpanel.Environment{
		Type:            environmentType,
		ConsoleTracker:  pufferpanel.CreateTracker(),
		StatusTracker:   pufferpanel.CreateTracker(),
		StatsTracker:    pufferpanel.CreateTracker(),
		ConsoleBuffer:   pufferpanel.CreateCache(),
		BackupDirectory: filepath.Join(backupFolder, server.Identifier),
		Wait:            &sync.Mutex{},
		Server:          server,
	}
	item.Implementation = factory()
	err := utils.UnmarshalTo(server.Environment.Metadata, item)
	if err != nil {
		return nil, err
	}

	err = utils.UnmarshalTo(server.Environment.Metadata, item.Implementation)
	if err != nil {
		return nil, err
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
