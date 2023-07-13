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
	"github.com/pufferpanel/pufferpanel/v3"
	"path/filepath"
	"sync"
)

var mapping = make(map[string]pufferpanel.EnvironmentFactory)

func Create(environmentType, folder, id string, environmentSection interface{}) (pufferpanel.Environment, error) {
	factory := mapping[environmentType]

	if factory == nil {
		switch environmentType {
		case "standard":
			fallthrough
		case "tty":
			factory = mapping["host"]
		}
	}

	if factory == nil {
		return nil, fmt.Errorf("undefined environment: %s", environmentType)
	}

	item := factory.Create(id)
	err := pufferpanel.UnmarshalTo(environmentSection, item)
	if err != nil {
		return nil, err
	}

	serverRoot := filepath.Join(folder, id)
	envCache := pufferpanel.CreateCache()
	wsManager := pufferpanel.CreateTracker()

	e := item.GetBase()
	if e.RootDirectory == "" {
		e.RootDirectory = serverRoot
	}
	e.WSManager = wsManager
	e.ConsoleBuffer = envCache
	e.Wait = &sync.WaitGroup{}
	e.Wrapper = e.CreateWrapper()

	return item, nil
}

func GetSupportedEnvironments() []string {
	result := make([]string, len(mapping))
	i := 0
	for k, _ := range mapping {
		result[i] = k
		i++
	}

	return result
}
