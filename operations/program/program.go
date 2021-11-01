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

package program

import (
	"fmt"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/programs"
	"time"
)

type Program struct {
	Action  string
	Program *programs.Program
}

func (d Program) Run(env pufferpanel.Environment) error {
	p := d.Program
	switch d.Action {
	case "install":
		return p.Install()
	case "start":
		return p.Start()
	case "stop":
		if err := p.Stop(); err != nil {
			return err
		}
		return WaitForStop(d)
	case "stopAsync":
		return p.Stop()
	case "restart":
		if err := p.Stop(); err != nil {
			return err
		}
		if err := WaitForStop(d); err != nil {
			return err
		}
		return p.Start()
	case "restartAsync":
		if err := p.Stop(); err != nil {
			return err
		}
		return p.Start()
	case "kill":
		return env.Kill()
	default:
		return fmt.Errorf("action %s was not valid action, expected one of: `install`, `start`,`stop`, `restart`, `kill`", d.Action)
	}
}

func WaitForStop(program Program) error {
	for {
		running, err := program.Program.IsRunning()
		if err != nil {
			return err
		}
		if !running {
			return nil
		}
		time.Sleep(time.Millisecond * 500)
	}
}
