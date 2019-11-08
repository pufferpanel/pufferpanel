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

package operations

import (
	"github.com/pufferpanel/pufferpanel/v2/daemon/programs/operations/ops"
	"github.com/pufferpanel/pufferpanel/v2/shared/logging"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
	"plugin"
	"reflect"
)

func loadOpModules() {
	var directory = path.Join(viper.GetString("data.modules"), "operations")

	files, err := ioutil.ReadDir(directory)
	if err != nil && os.IsNotExist(err) {
		return
	} else if err != nil {
		logging.Exception("error reading module directory", err)
	}

	for _, file := range files {
		logging.Info("Loading operation module: %s", file.Name())
		p, e := plugin.Open(path.Join(directory, file.Name()))
		if e != nil {
			logging.Exception("error opening module", err)
			continue
		}

		factory, e := p.Lookup("Factory")
		if e != nil {
			logging.Exception("error locating factory", err)
			continue
		}

		fty, ok := factory.(ops.OperationFactory)
		if !ok {
			logging.Error("expected OperationFactory, but found %s", reflect.TypeOf(factory).Name())
			continue
		}

		commandMapping[fty.Key()] = fty

		logging.Info("Loaded operation module: %s", fty.Key())
	}
}
