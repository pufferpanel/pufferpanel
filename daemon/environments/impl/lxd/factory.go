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

package lxd

import (
	"github.com/pufferpanel/pufferpanel/v2/shared"
	"github.com/pufferpanel/pufferpanel/v2/daemon/environments/envs"
	"github.com/pufferpanel/pufferpanel/v2/daemon/utils"
	"sync"
)

type EnvironmentFactory struct {
	envs.EnvironmentFactory
}

func (ef EnvironmentFactory) Create(folder, id string, environmentSection map[string]interface{}, rootDirectory string, cache shared.Cache, wsManager utils.WebSocketManager) envs.Environment {
	d := &lxc{BaseEnvironment: &envs.BaseEnvironment{Type: "lxc"}, TypeWithMetadata: shared.TypeWithMetadata{Type: "docker"}, ContainerId: id}
	d.BaseEnvironment.ExecutionFunction = d.executeAsync
	d.BaseEnvironment.WaitFunction = d.WaitForMainProcess
	d.wait = &sync.WaitGroup{}

	d.ImageName = shared.GetStringOrDefault(environmentSection, "image", "alpine/3.9")

	d.ContainerId = "pufferd-" + id
	d.RootDirectory = rootDirectory
	d.ConsoleBuffer = cache
	d.WSManager = wsManager
	return d
}

func (ef EnvironmentFactory) Key() string {
	return "lxd"
}
