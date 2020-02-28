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

package command

import (
	"fmt"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"strings"
)

type Command struct {
	Commands []string
	Env      map[string]string
}

func (c Command) Run(env pufferpanel.Environment) error {
	for _, cmd := range c.Commands {
		logging.Info().Printf("Executing command: %s", cmd)
		env.DisplayToConsole(true, fmt.Sprintf("Executing: %s\n", cmd))
		parts := strings.Split(cmd, " ")
		cmd := parts[0]
		args := parts[1:]
		_, err := env.Execute(cmd, args, c.Env, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
