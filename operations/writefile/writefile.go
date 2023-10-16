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

package writefile

import (
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"os"
	"path/filepath"
)

type WriteFile struct {
	TargetFile string
	Text       string
}

func (c WriteFile) Run(env pufferpanel.Environment) error {
	logging.Info.Printf("Writing data to file: %s", c.TargetFile)
	env.DisplayToConsole(true, "Writing some data to file: %s\n", c.TargetFile)
	target := filepath.Join(env.GetRootDirectory(), c.TargetFile)
	return os.WriteFile(target, []byte(c.Text), 0644)
}
