/*
 Copyright 2019 Padduck, LLC

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

package docker

import (
	"github.com/pufferpanel/pufferpanel/v2"
)

type EnvironmentFactory struct {
	pufferpanel.EnvironmentFactory
}

func (ef EnvironmentFactory) Create(id string) pufferpanel.Environment {
	d := &docker{
		BaseEnvironment: &pufferpanel.BaseEnvironment{Type: "docker", ServerId: id},
		ContainerId:     id,
		ImageName:       "pufferpanel/generic",
		NetworkMode:     "host",
		Ports:           make([]string, 0),
		Binds:           make(map[string]string),
	}

	d.ExecutionFunction = d.dockerExecuteAsync
	d.BaseEnvironment.WaitFunction = d.WaitForMainProcess
	return d
}

func (ef EnvironmentFactory) Key() string {
	return "docker"
}
