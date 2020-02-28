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

package programs

import (
	"container/list"
	"encoding/json"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/messages"
	"github.com/pufferpanel/pufferpanel/v2/operations"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Program struct {
	pufferpanel.Server

	CrashCounter int              `json:"-"`
	Environment  pufferpanel.Environment `json:"-"`
}

var queue *list.List
var lock = sync.Mutex{}
var ticker *time.Ticker
var running = false

func InitService() {
	queue = list.New()
	ticker = time.NewTicker(1 * time.Second)
	running = true
	go processQueue()
}

func StartViaService(p *Program) {
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()

	if running {
		queue.PushBack(p)
	}
}

func ShutdownService() {
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()

	running = false
	ticker.Stop()
}

func processQueue() {
	for range ticker.C {
		lock.Lock()
		next := queue.Front()
		if next != nil {
			queue.Remove(next)
		}
		lock.Unlock()
		if next == nil {
			continue
		}
		program := next.Value.(*Program)
		if run, _ := program.IsRunning(); !run {
			err := program.Start()
			if err != nil {
				logging.Error().Printf("Error starting server %s: %s", program.Id(), err)
			}
		}
	}
}

type FileData struct {
	Contents      io.ReadCloser
	ContentLength int64
	FileList      []messages.FileDesc
	Name          string
}

func (p *Program) DataToMap() map[string]interface{} {
	var result = make(map[string]interface{}, len(p.Variables))

	for k, v := range p.Variables {
		result[k] = v.Value
	}

	return result
}

func CreateProgram() *Program {
	return &Program{
		Server: pufferpanel.Server{
			Execution: pufferpanel.Execution{
				Disabled:                false,
				AutoStart:               false,
				AutoRestartFromCrash:    false,
				AutoRestartFromGraceful: false,
				PreExecution:            make([]interface{}, 0),
				PostExecution:           make([]interface{}, 0),
				EnvironmentVariables:    make(map[string]string, 0),
			},
			Type:           "standard",
			Variables:      make(map[string]pufferpanel.Variable, 0),
			Display:        "Unknown server",
			Installation:   make([]interface{}, 0),
			Uninstallation: make([]interface{}, 0),
		},
	}
}

//Starts the program.
//This includes starting the environment if it is not running.
func (p *Program) Start() (err error) {
	if !p.IsEnabled() {
		logging.Error().Printf("Server %s is not enabled, cannot start", p.Id())
		return pufferpanel.ErrServerDisabled
	}
	if running, err := p.IsRunning(); running || err != nil {
		return err
	}

	logging.Info().Printf("Starting server %s", p.Id())
	p.Environment.DisplayToConsole(true, "Starting server\n")
	data := make(map[string]interface{})
	for k, v := range p.Variables {
		data[k] = v.Value
	}

	process, err := operations.GenerateProcess(p.Execution.PreExecution, p.Environment, p.DataToMap(), p.Execution.EnvironmentVariables)
	if err != nil {
		logging.Error().Printf("Error generating pre-execution steps: %s", err)
		p.Environment.DisplayToConsole(true, "Error running pre execute\n")
		p.Environment.DisplayToConsole(true, "%s\n", err.Error())
		return
	}

	err = process.Run(p.Environment)
	if err != nil {
		logging.Error().Printf("Error running pre-execution steps: %s", err)
		p.Environment.DisplayToConsole(true, "Error running pre execute\n")
		p.Environment.DisplayToConsole(true, "%s\n", err.Error())
		return
	}

	//HACK: add rootDir stuff
	data["rootDir"] = p.Environment.GetRootDirectory()

	err = p.Environment.ExecuteAsync(p.Execution.ProgramName, pufferpanel.ReplaceTokensInArr(p.Execution.Arguments, data), pufferpanel.ReplaceTokensInMap(p.Execution.EnvironmentVariables, data), p.afterExit)
	if err != nil {
		logging.Error().Printf("error starting server %s: %s", p.Id(), err)
		p.Environment.DisplayToConsole(true, " Failed to start server\n")
		p.Environment.DisplayToConsole(true, "%s\n", err.Error())
	}

	return
}

//Stops the program.
//This will also stop the environment it is ran in.
func (p *Program) Stop() (err error) {
	if running, err := p.IsRunning(); !running || err != nil {
		return err
	}

	logging.Info().Printf("Stopping server %s", p.Id())
	if p.Execution.StopCode != 0 {
		err = p.Environment.SendCode(p.Execution.StopCode)
	} else {
		err = p.Environment.ExecuteInMainProcess(p.Execution.StopCommand)
	}
	if err != nil {
		logging.Error().Printf("Error stopping server: %s", err)
		p.Environment.DisplayToConsole(true, "Failed to stop server\n")
		p.Environment.DisplayToConsole(true, "%s\n", err.Error())
	} else {
		p.Environment.DisplayToConsole(true, "Server was told to stop\n")
	}
	return
}

//Kills the program.
//This will also stop the environment it is ran in.
func (p *Program) Kill() (err error) {
	logging.Info().Printf("Killing server %s", p.Id())
	err = p.Environment.Kill()
	if err != nil {
		logging.Error().Printf("Error killing server: %s", err)
		p.Environment.DisplayToConsole(true, "Failed to kill server\n")
		p.Environment.DisplayToConsole(true, "%s\n", err.Error())
	} else {
		p.Environment.DisplayToConsole(true, "Server killed\n")
	}
	return
}

//Creates any files needed for the program.
//This includes creating the environment.
func (p *Program) Create() (err error) {
	logging.Info().Printf("Creating server %s", p.Id())
	p.Environment.DisplayToConsole(true, "Allocating server\n")
	err = p.Environment.Create()
	if err != nil {
		logging.Error().Printf("Error creating server: %s", err)
		p.Environment.DisplayToConsole(true, "Failed to create server\n")
		p.Environment.DisplayToConsole(true, "%s\n", err.Error())
	} else {
		p.Environment.DisplayToConsole(true, "Server allocated\n")
		p.Environment.DisplayToConsole(true, "Ready to be installed\n")
	}

	return
}

//Destroys the server.
//This will delete the server, environment, and any files related to it.
func (p *Program) Destroy() (err error) {
	logging.Info().Printf("Destroying server %s", p.Id())
	process, err := operations.GenerateProcess(p.Uninstallation, p.Environment, p.DataToMap(), p.Execution.EnvironmentVariables)
	if err != nil {
		logging.Error().Printf("Error uninstalling server: %s", err)
		p.Environment.DisplayToConsole(true, "Failed to uninstall server\n")
		p.Environment.DisplayToConsole(true, "%s\n", err.Error())
		return
	}

	err = process.Run(p.Environment)
	if err != nil {
		logging.Error().Printf("Error uninstalling server: %s", err)
		p.Environment.DisplayToConsole(true, "Failed to uninstall server\n")
		p.Environment.DisplayToConsole(true, "%s\n", err.Error())
		return
	}

	err = p.Environment.Delete()
	if err != nil {
		logging.Error().Printf("Error uninstalling server: %s", err)
		p.Environment.DisplayToConsole(true, "Failed to uninstall server\n")
		p.Environment.DisplayToConsole(true, "%s\n", err.Error())
	}
	return
}

func (p *Program) Install() (err error) {
	if !p.IsEnabled() {
		logging.Error().Printf("Server %s is not enabled, cannot install", p.Id())
		return pufferpanel.ErrServerDisabled
	}

	logging.Info().Printf("Installing server %s", p.Id())
	running, err := p.IsRunning()
	if err != nil {
		logging.Error().Printf("Error checking server status: %s", err)
		p.Environment.DisplayToConsole(true, "Error on checking to see if server is running\n")
		return
	}

	if running {
		err = p.Stop()
	}

	if err != nil {
		logging.Error().Printf("Error stopping server: %s", err)
		p.Environment.DisplayToConsole(true, "Failed to stop server\n")
		p.Environment.DisplayToConsole(true, "%s\n", err.Error())
		return
	}

	p.Environment.DisplayToConsole(true, "Installing server\n")

	err = os.MkdirAll(p.Environment.GetRootDirectory(), 0755)
	if err != nil && !os.IsExist(err) {
		logging.Error().Printf("Error creating server directory: %s", err)
		p.Environment.DisplayToConsole(true, "Failed to create server directory\n")
		p.Environment.DisplayToConsole(true, "%s\n", err.Error())
		return
	}

	if len(p.Installation) > 0 {
		process, err := operations.GenerateProcess(p.Installation, p.GetEnvironment(), p.DataToMap(), p.Execution.EnvironmentVariables)
		if err != nil {
			logging.Error().Printf("Error installing server: %s", err)
			p.Environment.DisplayToConsole(true, "Failed to install server\n")
			p.Environment.DisplayToConsole(true, "%s\n", err.Error())
			return err
		}

		err = process.Run(p.Environment)
		if err != nil {
			logging.Error().Printf("Error installing server: %s", err)
			p.Environment.DisplayToConsole(true, "Failed to install server\n")
			p.Environment.DisplayToConsole(true, "%s\n", err.Error())
			return err
		}
	}

	p.Environment.DisplayToConsole(true, "Server installed\n")
	return
}

//Determines if the server is running.
func (p *Program) IsRunning() (isRunning bool, err error) {
	isRunning, err = p.Environment.IsRunning()
	return
}

//Sends a command to the process
//If the program supports input, this will send the arguments to that.
func (p *Program) Execute(command string) (err error) {
	err = p.Environment.ExecuteInMainProcess(command)
	return
}

func (p *Program) SetEnabled(isEnabled bool) (err error) {
	p.Execution.Disabled = !isEnabled
	return
}

func (p *Program) IsEnabled() (isEnabled bool) {
	return !p.Execution.Disabled
}

func (p *Program) SetEnvironment(environment pufferpanel.Environment) (err error) {
	p.Environment = environment
	return
}

func (p *Program) Id() string {
	return p.Identifier
}

func (p *Program) GetEnvironment() pufferpanel.Environment {
	return p.Environment
}

func (p *Program) SetAutoStart(isAutoStart bool) (err error) {
	p.Execution.AutoStart = isAutoStart
	return
}

func (p *Program) IsAutoStart() (isAutoStart bool) {
	isAutoStart = p.Execution.AutoStart
	return
}

func (p *Program) Save(file string) (err error) {
	logging.Info().Printf("Saving server %s", p.Id())

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(file, data, 0664)
	return
}

func (p *Program) Edit(data map[string]pufferpanel.Variable, overrideUser bool) (err error) {
	for k, v := range data {
		var elem pufferpanel.Variable

		if _, ok := p.Variables[k]; ok {
			elem = p.Variables[k]
		} else {
			//copy from provided
			elem = v
		}
		if !elem.UserEditable && !overrideUser {
			continue
		}

		elem.Value = v.Value

		p.Variables[k] = elem
	}
	err = Save(p.Id())
	return
}

func (p *Program) GetData() map[string]pufferpanel.Variable {
	return p.Variables
}

func (p *Program) GetNetwork() string {
	data := p.GetData()
	ip := "0.0.0.0"
	port := "0"

	if ipData, ok := data["ip"]; ok {
		if _, ok = ipData.Value.(string); ok {
			ip = ipData.Value.(string)
		}
	}

	if portData, ok := data["port"]; ok {
		if _, ok = portData.Value.(string); ok {
			port = portData.Value.(string)
		}
	}

	return ip + ":" + port
}

func (p *Program) CopyFrom(s *Program) {
	p.Variables = s.Variables
	p.Execution = s.Execution
	p.Display = s.Display
	p.Environment = s.Environment
	p.Installation = s.Installation
	p.Uninstallation = s.Uninstallation
	p.Type = s.Type
}

func (p *Program) afterExit(graceful bool) {
	if graceful {
		p.CrashCounter = 0
	}

	mapping := p.DataToMap()
	mapping["success"] = graceful

	processes, err := operations.GenerateProcess(p.Execution.PostExecution, p.Environment, mapping, p.Execution.EnvironmentVariables)
	if err != nil {
		logging.Error().Printf("Error running post processing for server %s: %s", p.Id(), err)
		p.Environment.DisplayToConsole(true, "Failed to run post-execution steps\n")
		p.Environment.DisplayToConsole(true, "%s\n", err.Error())
		return
	}
	p.Environment.DisplayToConsole(true, "Running post-execution steps\n")
	logging.Info().Printf("Running post execution steps: %s", p.Id())

	err = processes.Run(p.Environment)
	if err != nil {
		logging.Error().Printf("Error running post processing for server: %s", err)
		p.Environment.DisplayToConsole(true, "Failed to run post-execution steps\n")
		p.Environment.DisplayToConsole(true, "%s\n", err.Error())
		return
	}

	if !p.Execution.AutoRestartFromCrash && !p.Execution.AutoRestartFromGraceful {
		return
	}

	if graceful && p.Execution.AutoRestartFromGraceful {
		StartViaService(p)
	} else if !graceful && p.Execution.AutoRestartFromCrash && p.CrashCounter < viper.GetInt("daemon.data.crashLimit") {
		p.CrashCounter++
		StartViaService(p)
	}
}

func (p *Program) GetItem(name string) (*FileData, error) {
	targetFile := pufferpanel.JoinPath(p.GetEnvironment().GetRootDirectory(), name)
	if !pufferpanel.EnsureAccess(targetFile, p.GetEnvironment().GetRootDirectory()) {
		return nil, pufferpanel.ErrIllegalFileAccess
	}

	info, err := os.Stat(targetFile)

	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		files, _ := ioutil.ReadDir(targetFile)
		var fileNames []messages.FileDesc
		offset := 0
		if name == "" || name == "." || name == "/" {
			fileNames = make([]messages.FileDesc, len(files))
		} else {
			fileNames = make([]messages.FileDesc, len(files)+1)
			fileNames[0] = messages.FileDesc{
				Name: "..",
				File: false,
			}
			offset = 1
		}

		//validate any symlinks are valid
		files = pufferpanel.RemoveInvalidSymlinks(files, targetFile, p.GetEnvironment().GetRootDirectory())

		for i, file := range files {
			newFile := messages.FileDesc{
				Name: file.Name(),
				File: !file.IsDir(),
			}

			if newFile.File {
				newFile.Size = file.Size()
				newFile.Modified = file.ModTime().Unix()
				newFile.Extension = filepath.Ext(file.Name())
			}

			fileNames[i+offset] = newFile
		}

		return &FileData{FileList: fileNames}, nil
	} else {
		file, err := os.Open(targetFile)
		if err != nil {
			return nil, err
		}
		return &FileData{Contents: file, ContentLength: info.Size(), Name: info.Name()}, nil
	}
}

func (p *Program) CreateFolder(name string) error {
	folder := pufferpanel.JoinPath(p.GetEnvironment().GetRootDirectory(), name)

	if !pufferpanel.EnsureAccess(folder, p.GetEnvironment().GetRootDirectory()) {
		return pufferpanel.ErrIllegalFileAccess
	}
	return os.Mkdir(folder, 0755)
}

func (p *Program) OpenFile(name string) (io.WriteCloser, error) {
	targetFile := pufferpanel.JoinPath(p.GetEnvironment().GetRootDirectory(), name)

	if !pufferpanel.EnsureAccess(targetFile, p.GetEnvironment().GetRootDirectory()) {
		return nil, pufferpanel.ErrIllegalFileAccess
	}

	file, err := os.Create(targetFile)
	return file, err
}

func (p *Program) DeleteItem(name string) error {
	targetFile := pufferpanel.JoinPath(p.GetEnvironment().GetRootDirectory(), name)

	if !pufferpanel.EnsureAccess(targetFile, p.GetEnvironment().GetRootDirectory()) {
		return pufferpanel.ErrIllegalFileAccess
	}

	return os.RemoveAll(targetFile)
}
