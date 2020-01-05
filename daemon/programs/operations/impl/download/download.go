/*
 Copyright 2016 Padduck, LLC

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

package download

import (
	"github.com/cavaliercoder/grab"
	"github.com/pufferpanel/pufferpanel/v2/daemon/environments/envs"
	"github.com/pufferpanel/pufferpanel/v2/logging"
)

type Download struct {
	Files []string
}

func (d Download) Run(env envs.Environment) error {
	for _, file := range d.Files {
		logging.Info().Printf("Download file from %s to %s", file, env.GetRootDirectory())
		env.DisplayToConsole(true, "Downloading file %s\n", file)
		_, err := grab.Get(env.GetRootDirectory(), file)
		if err != nil {
			return err
		}
	}
	return nil
}
