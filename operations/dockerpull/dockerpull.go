/*
 Copyright 2023 PufferPanel

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

package dockerpull

import (
	"context"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/servers/docker"
)

type DockerPull struct {
	ImageName string
}

func (d DockerPull) Run(env pufferpanel.Environment) error {
	dockerEnv, ok := env.(*docker.Docker)

	if !ok {
		return pufferpanel.ErrEnvironmentNotSupported
	}

	return dockerEnv.PullImage(context.Background(), d.ImageName, true)
}
