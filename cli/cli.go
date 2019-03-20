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

package cli

import (
	"github.com/pufferpanel/apufferi/cli"
	"github.com/pufferpanel/pufferpanel/cli/commands"
)

func Run() error {
	//define command order here
	cmds := []cli.Command{
		&commands.Run{},
	}

	return cli.Run(cmds)
}
