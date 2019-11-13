/*
 Copyright 2018 Padduck, LLC

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 	http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package environments

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/daemon/cache"
	"github.com/pufferpanel/pufferpanel/v2/daemon/environments/envs"
	"github.com/pufferpanel/pufferpanel/v2/daemon/environments/impl/docker"
	"github.com/pufferpanel/pufferpanel/v2/daemon/environments/impl/standard"
	"github.com/pufferpanel/pufferpanel/v2/daemon/utils"
	"path/filepath"
	"sync"
)

var mapping map[string]envs.EnvironmentFactory

func LoadModules() {
	mapping = make(map[string]envs.EnvironmentFactory)

	mapping["standard"] = standard.EnvironmentFactory{}
	mapping["docker"] = docker.EnvironmentFactory{}

	loadAdditionalModules(mapping)
}

func Create(environmentType, folder, id string, environmentSection interface{}) (envs.Environment, error) {
	factory := mapping[environmentType]

	if factory == nil {
		return nil, errors.New(fmt.Sprintf("undefined environment: %s", environmentType))
	}

	item := factory.Create(id)
	err := pufferpanel.UnmarshalTo(environmentSection, item)
	if err != nil {
		return nil, err
	}

	serverRoot := filepath.Join(folder, id)
	envCache := cache.CreateCache()
	wsManager := utils.CreateWSManager()

	e := item.GetBase()
	if e.RootDirectory == "" {
		e.RootDirectory = serverRoot
	}
	e.WSManager = wsManager
	e.ConsoleBuffer = envCache
	e.Wait = &sync.WaitGroup{}

	return item, nil
}

func GetSupportedEnvironments() []string {
	result := make([]string, len(mapping))
	i := 0
	for _, v := range mapping {
		result[i] = v.Key()
		i++
	}

	return result
}
