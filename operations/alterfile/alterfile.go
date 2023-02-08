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

package alterfile

import (
	"bytes"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

type AlterFile struct {
	TargetFile string
	Search     string
	Replace    string
	Regex      bool
}

func (c AlterFile) Run(env pufferpanel.Environment) error {
	logging.Info.Printf("Changing data in file: %s", c.TargetFile)
	env.DisplayToConsole(true, "Changing some data in file: %s\n ", c.TargetFile)
	target := filepath.Join(env.GetRootDirectory(), c.TargetFile)
	data, err := ioutil.ReadFile(target)
	if err != nil {
		return err
	}

	var out []byte
	if c.Regex {
		regex, err := regexp.Compile("(?m)" + c.Search)
		if err != nil {
			return err
		}
		out = regex.ReplaceAllLiteral(data, []byte(c.Replace))
	} else {
		out = bytes.ReplaceAll(data, []byte(c.Search), []byte(c.Replace))
	}

	return ioutil.WriteFile(target, out, 0644)
}
