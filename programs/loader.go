/*
 Copyright 2016 Padduck, LLC

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 distributed under the License is distributed on an "AS IS" BASIS,
 You may obtain a copy of the License at

 	http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package programs

import (
	"encoding/json"
	"errors"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/environments"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/operations"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	allPrograms = make([]*Program, 0)
)

func Initialize() {
	operations.LoadOperations()
}

func LoadFromFolder() {
	err := os.Mkdir(pufferpanel.ServerFolder, 0755)
	if err != nil && !os.IsExist(err) {
		logging.Error().Fatalf("Error creating server data folder: %s", err)
	}
	programFiles, err := ioutil.ReadDir(pufferpanel.ServerFolder)
	if err != nil {
		logging.Error().Fatalf("Error reading from server data folder: %s", err)
	}
	var program *Program
	for _, element := range programFiles {
		if element.IsDir() {
			continue
		}
		logging.Info().Printf("Attempting to load " + element.Name())
		id := strings.TrimSuffix(element.Name(), filepath.Ext(element.Name()))
		program, err = Load(id)
		if err != nil {
			logging.Error().Printf("Error loading server details from json (%s): %s", element.Name(), err)
			continue
		}
		program.Scheduler = NewScheduler(program)
		err = program.Scheduler.LoadMap(program.Tasks)
		if err != nil {
			logging.Error().Printf("Error loading server tasks from json (%s): %s", element.Name(), err)
			continue
		}
		err = program.Scheduler.Start()
		if err != nil {
			logging.Error().Printf("Error starting server scheduler (%s): %s", element.Name(), err)
			continue
		}
		logging.Info().Printf("Loaded server %s", program.Id())
		allPrograms = append(allPrograms, program)
	}
}

func Get(id string) (program *Program, err error) {
	program = GetFromCache(id)
	if program == nil {
		program, err = Load(id)
	}
	return
}

func GetAll() []*Program {
	return allPrograms
}

func Load(id string) (program *Program, err error) {
	var data []byte
	data, err = ioutil.ReadFile(filepath.Join(pufferpanel.ServerFolder, id+".json"))
	if len(data) == 0 || err != nil {
		return
	}
	program, err = LoadFromData(id, data)
	return
}

func LoadFromData(id string, source []byte) (*Program, error) {
	data := CreateProgram()

	//HACK: Because golang thinks environment and Environment in the json are the same, we have to manually clean the
	//invalid record up....
	rawMap := make(map[string]interface{})
	err := json.Unmarshal(source, &rawMap)
	if err != nil {
		return nil, err
	}

	delete(rawMap, "Environment")
	source, err = json.Marshal(rawMap)

	err = json.Unmarshal(source, &data)
	if err != nil {
		return nil, err
	}

	data.Identifier = id

	if data.Execution.LegacyRun != "" {
		data.Execution.Command = strings.TrimSpace(data.Execution.LegacyRun + " " + strings.Join(data.Execution.LegacyArguments, " "))
		data.Execution.LegacyRun = ""
		data.Execution.LegacyArguments = nil
		err = data.Save()
		if err != nil {
			return nil, err
		}
	}

	var typeMap pufferpanel.Type
	err = pufferpanel.UnmarshalTo(data.Environment, &typeMap)
	if err != nil {
		return nil, err
	}

	environmentType := typeMap.Type
	data.RunningEnvironment, err = environments.Create(environmentType, pufferpanel.ServerFolder, id, data.Environment)
	return data, nil
}

func Create(program *Program) error {
	if GetFromCache(program.Id()) != nil {
		return pufferpanel.ErrServerAlreadyExists
	}

	var err error

	defer func() {
		if err != nil {
			//revert since we have an error
			_ = os.Remove(filepath.Join(pufferpanel.ServerFolder, program.Id()+".json"))
			if program.RunningEnvironment != nil {
				_ = program.RunningEnvironment.Delete()
			}
		}
	}()

	f, err := os.Create(filepath.Join(pufferpanel.ServerFolder, program.Id()+".json"))
	defer pufferpanel.Close(f)
	if err != nil {
		logging.Error().Printf("Error writing server: %s", err)
		return err
	}

	encoder := json.NewEncoder(f)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(program)

	if err != nil {
		logging.Error().Printf("Error writing server: %s", err)
		return err
	}

	var typeMap pufferpanel.Type
	err = pufferpanel.UnmarshalTo(program.Environment, &typeMap)
	if err != nil {
		return err
	}

	program.RunningEnvironment, err = environments.Create(typeMap.Type, pufferpanel.ServerFolder, program.Id(), program.Environment)

	err = program.Create()
	if err != nil {
		return err
	}

	allPrograms = append(allPrograms, program)
	return nil
}

func Delete(id string) (err error) {
	var index int
	var program *Program
	for i, element := range allPrograms {
		if element.Id() == id {
			program = element
			index = i
			break
		}
	}
	if program == nil {
		return
	}
	running, err := program.IsRunning()

	if err != nil {
		return
	}

	if running {
		err = program.Kill()
		if err != nil {
			return
		}

		err = program.RunningEnvironment.WaitForMainProcess()
		if err != nil {
			return
		}
	}

	_ = program.Scheduler.Stop()

	err = program.Destroy()
	if err != nil {
		return
	}
	err = os.Remove(filepath.Join(pufferpanel.ServerFolder, program.Id()+".json"))
	if err != nil {
		logging.Error().Printf("Error removing server: %s", err)
	}
	allPrograms = append(allPrograms[:index], allPrograms[index+1:]...)
	return
}

func GetFromCache(id string) *Program {
	for _, element := range allPrograms {
		if element != nil && element.Id() == id {
			return element
		}
	}
	return nil
}

func Save(id string) (err error) {
	program := GetFromCache(id)
	if program == nil {
		err = errors.New("no server with given id")
		return
	}
	err = program.Save()
	return
}

func Reload(id string) (err error) {
	program := GetFromCache(id)
	if program == nil {
		err = errors.New("server does not exist")
		return
	}

	logging.Info().Printf("Reloading server %s", program.Id())
	newVersion, err := Load(id)
	if err != nil {
		logging.Error().Printf("Error reloading server: %s", err)
		return
	}

	program.RunningEnvironment = newVersion.RunningEnvironment
	program.Server = newVersion.Server

	_ = program.Scheduler.Stop()
	logging.Debug().Println("Rebuilding scheduler")
	err = newVersion.Scheduler.Rebuild()
	if err != nil {
		logging.Error().Printf("Error reloading server scheduler: %s", err)
		return err
	}

	logging.Debug().Println("Loading scheduled tasks")
	err = newVersion.Scheduler.LoadMap(program.Tasks)
	if err != nil {
		return err
	}

	logging.Debug().Println("Starting scheduler")
	err = newVersion.Scheduler.Start()
	if err != nil {
		return err
	}

	return
}
