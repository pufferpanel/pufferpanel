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

package servers

import (
	"encoding/json"
	"errors"
	"github.com/go-co-op/gocron"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	allServers = make([]*Server, 0)
)

func LoadFromFolder() {
	err := os.Mkdir(config.ServersFolder.Value(), 0755)
	if err != nil && !os.IsExist(err) {
		logging.Error.Fatalf("Error creating server data folder: %s", err)
	}
	programFiles, err := os.ReadDir(config.ServersFolder.Value())
	if err != nil {
		logging.Error.Fatalf("Error reading from server data folder: %s", err)
	}
	var program *Server
	for _, element := range programFiles {
		if element.IsDir() || !strings.HasSuffix(element.Name(), ".json") {
			continue
		}
		logging.Info.Printf("Attempting to load " + element.Name())
		id := strings.TrimSuffix(element.Name(), filepath.Ext(element.Name()))
		program, err = Load(id)
		if err != nil {
			logging.Error.Printf("Error loading server details from json (%s): %s", element.Name(), err)
			continue
		}

		err = program.Scheduler.LoadMap(program.Tasks)
		if err != nil {
			logging.Error.Printf("Error loading server tasks from json (%s): %s", element.Name(), err)
			continue
		}
		err = program.Scheduler.Start()
		if err != nil {
			logging.Error.Printf("Error starting server scheduler (%s): %s", element.Name(), err)
			continue
		}
		logging.Info.Printf("Loaded server %s", program.Id())
		allServers = append(allServers, program)
	}
}

func Get(id string) (program *Server, err error) {
	program = GetFromCache(id)
	if program == nil {
		program, err = Load(id)
	}
	return
}

func GetAll() []*Server {
	return allServers
}

func Load(id string) (program *Server, err error) {
	var data []byte
	data, err = os.ReadFile(filepath.Join(config.ServersFolder.Value(), id+".json"))
	if len(data) == 0 || err != nil {
		return
	}
	program, err = LoadFromData(id, data)
	return
}

func LoadFromData(id string, source []byte) (*Server, error) {
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
	if err != nil {
		return nil, err
	}

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

	if environmentType == "standard" || environmentType == "tty" {
		environmentType = "host"

		tracker := make(map[string]interface{})
		err = pufferpanel.UnmarshalTo(data.Environment, &tracker)
		if err != nil {
			return nil, err
		}
		tracker["type"] = "host"
		data.Environment = tracker

		err = data.Save()
		if err != nil {
			return nil, err
		}
	}

	data.RunningEnvironment, err = CreateEnvironment(environmentType, config.ServersFolder.Value(), id, data.Environment)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func startScheduler(program *Server) error {
	s := gocron.NewScheduler(time.UTC)
	for k, t := range program.Tasks {
		if t.CronSchedule != "" {
			_, err := s.Cron(t.CronSchedule).Do(executeTask, program, k)
			if err != nil {
				return err
			}
		}
	}
	s.SetMaxConcurrentJobs(5, gocron.RescheduleMode)
	s.StartAsync()
	program.Scheduler.scheduler = s
	return nil
}

func executeTask(p *Server, taskId string) (err error) {
	task, ok := p.Tasks[taskId]
	if !ok {
		logging.Error.Printf("Execution of task %s on server %s requested, but task not found", taskId, p.Id())
		return
	}

	ops := task.Operations
	if len(ops) > 0 {
		p.RunningEnvironment.DisplayToConsole(true, "Running task %s\n", task.Name)
		var process OperationProcess
		process, err = GenerateProcess(ops, p.GetEnvironment(), p.DataToMap(), p.Execution.EnvironmentVariables)
		if err != nil {
			logging.Error.Printf("Error setting up tasks: %s", err)
			p.RunningEnvironment.DisplayToConsole(true, "Failed to setup tasks\n")
			p.RunningEnvironment.DisplayToConsole(true, "%s\n", err.Error())
			return
		}

		err = process.Run(p)
		if err != nil {
			logging.Error.Printf("Error setting up tasks: %s", err)
			p.RunningEnvironment.DisplayToConsole(true, "Failed to setup tasks\n")
			p.RunningEnvironment.DisplayToConsole(true, "%s\n", err.Error())
			return
		}
		p.RunningEnvironment.DisplayToConsole(true, "Task %s finished\n", task.Name)
	}
	return
}

func ExecuteTask(programId string, taskId string) error {
	program := GetFromCache(programId)
	if program == nil {
		return errors.New("no server with given id")
	}
	return executeTask(program, taskId)
}

func Create(program *Server) error {
	if GetFromCache(program.Id()) != nil {
		return pufferpanel.ErrServerAlreadyExists
	}

	var err error

	defer func() {
		if err != nil {
			//revert since we have an error
			_ = os.Remove(filepath.Join(config.ServersFolder.Value(), program.Id()+".json"))
			if program.RunningEnvironment != nil {
				_ = program.RunningEnvironment.Delete()
			}
		}
	}()

	f, err := os.Create(filepath.Join(config.ServersFolder.Value(), program.Id()+".json"))
	defer pufferpanel.Close(f)
	if err != nil {
		logging.Error.Printf("Error writing server: %s", err)
		return err
	}

	encoder := json.NewEncoder(f)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(program)

	if err != nil {
		logging.Error.Printf("Error writing server: %s", err)
		return err
	}

	var typeMap pufferpanel.Type
	err = pufferpanel.UnmarshalTo(program.Environment, &typeMap)
	if err != nil {
		return err
	}

	program.RunningEnvironment, err = CreateEnvironment(typeMap.Type, config.ServersFolder.Value(), program.Id(), program.Environment)
	if err != nil {
		return err
	}

	err = program.Create()
	if err != nil {
		return err
	}

	allServers = append(allServers, program)
	return nil
}

func Delete(id string) (err error) {
	var index int
	var program *Server
	for i, element := range allServers {
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
	err = os.Remove(filepath.Join(config.ServersFolder.Value(), program.Id()+".json"))
	if err != nil {
		logging.Error.Printf("Error removing server: %s", err)
	}
	allServers = append(allServers[:index], allServers[index+1:]...)
	return
}

func GetFromCache(id string) *Server {
	for _, element := range allServers {
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

	logging.Info.Printf("Reloading server %s", program.Id())
	newVersion, err := Load(id)
	if err != nil {
		logging.Error.Printf("Error reloading server: %s", err)
		return
	}

	program.RunningEnvironment = newVersion.RunningEnvironment
	program.Server = newVersion.Server

	_ = program.Scheduler.Stop()
	logging.Debug.Println("Rebuilding scheduler")
	err = newVersion.Scheduler.Rebuild()
	if err != nil {
		logging.Error.Printf("Error reloading server scheduler: %s", err)
		return err
	}

	logging.Debug.Println("Loading scheduled tasks")
	err = newVersion.Scheduler.LoadMap(program.Tasks)
	if err != nil {
		return err
	}

	logging.Debug.Println("Starting scheduler")
	err = newVersion.Scheduler.Start()
	if err != nil {
		return err
	}

	return
}
