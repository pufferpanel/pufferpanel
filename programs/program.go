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
	"fmt"
	"github.com/mholt/archiver/v3"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/config"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/messages"
	"github.com/pufferpanel/pufferpanel/v2/operations"
	"github.com/spf13/cast"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

type Program struct {
	pufferpanel.Server

	CrashCounter       int                     `json:"-"`
	RunningEnvironment pufferpanel.Environment `json:"-"`
	Scheduler          Scheduler               `json:"-"`
	fs                 pufferpanel.FileServer  `json:"-"`
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
	if !running {
		return
	}

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
				logging.Error.Printf("[%s] Error starting server: %s", program.Id(), err)
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
	var result = p.Server.DataToMap()
	result["rootDir"] = p.RunningEnvironment.GetRootDirectory()

	return result
}

func CreateProgram() *Program {
	p := &Program{
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
			Type:           pufferpanel.Type{Type: "standard"},
			Variables:      make(map[string]pufferpanel.Variable, 0),
			Tasks:          make(map[string]pufferpanel.Task, 0),
			Display:        "Unknown server",
			Installation:   make([]interface{}, 0),
			Uninstallation: make([]interface{}, 0),
		},
	}
	p.Scheduler = NewScheduler(p)
	return p
}

// Starts the program.
// This includes starting the environment if it is not running.
func (p *Program) Start() (err error) {
	if !p.IsEnabled() {
		p.Log(logging.Error, "Server %s is not enabled, cannot start", p.Id())
		return pufferpanel.ErrServerDisabled
	}
	if running, err := p.IsRunning(); running || err != nil {
		return err
	}

	p.Log(logging.Info, "Starting server %s", p.Id())
	p.RunningEnvironment.DisplayToConsole(true, "Starting server\n")

	data := p.DataToMap()

	process, err := operations.GenerateProcess(p.Execution.PreExecution, p.RunningEnvironment, data, p.Execution.EnvironmentVariables)
	if err != nil {
		p.Log(logging.Error, "Error generating pre-execution steps: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Error running pre execute\n")
		return
	}

	err = process.Run(p.RunningEnvironment)
	if err != nil {
		p.Log(logging.Error, "Error running pre-execution steps: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Error running pre execute\n")
		return
	}

	commandLine := pufferpanel.ReplaceTokens(p.Execution.Command, data)
	if p.Execution.WorkingDirectory == "${rootDir}" {
		p.Execution.WorkingDirectory = ""
	}
	workDir := pufferpanel.ReplaceTokens(p.Execution.WorkingDirectory, data)

	if !pufferpanel.EnsureAccess(path.Join(p.RunningEnvironment.GetRootDirectory(), workDir), p.RunningEnvironment.GetRootDirectory()) {
		p.Log(logging.Error, "Working directory is invalid for server: %s", workDir)
		p.RunningEnvironment.DisplayToConsole(true, "Working directory is invalid for server: %s", workDir)
		return
	}

	cmd, args := pufferpanel.SplitArguments(commandLine)
	err = p.RunningEnvironment.ExecuteAsync(pufferpanel.ExecutionData{
		Command:          cmd,
		Arguments:        args,
		Environment:      pufferpanel.ReplaceTokensInMap(p.Execution.EnvironmentVariables, data),
		WorkingDirectory: workDir,
		Callback:         p.afterExit,
	})

	if err != nil {
		p.Log(logging.Error, "error starting server %s: %s", p.Id(), err)
		p.RunningEnvironment.DisplayToConsole(true, " Failed to start server\n")
	}

	return
}

// Stops the program.
// This will also stop the environment it is ran in.
func (p *Program) Stop() (err error) {
	if running, err := p.IsRunning(); !running || err != nil {
		return err
	}

	p.Log(logging.Info, "Stopping server %s", p.Id())
	if p.Execution.StopCode != 0 {
		err = p.RunningEnvironment.SendCode(p.Execution.StopCode)
	} else {
		err = p.RunningEnvironment.ExecuteInMainProcess(p.Execution.StopCommand)
	}
	if err != nil {
		p.Log(logging.Error, "Error stopping server: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Failed to stop server\n")
	} else {
		p.RunningEnvironment.DisplayToConsole(true, "Server was told to stop\n")
	}
	return
}

// Kills the program.
// This will also stop the environment it is ran in.
func (p *Program) Kill() (err error) {
	p.Log(logging.Info, "Killing server %s", p.Id())
	err = p.RunningEnvironment.Kill()
	if err != nil {
		p.Log(logging.Error, "Error killing server: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Failed to kill server\n")
	} else {
		p.RunningEnvironment.DisplayToConsole(true, "Server killed\n")
	}
	return
}

// Creates any files needed for the program.
// This includes creating the environment.
func (p *Program) Create() (err error) {
	p.Log(logging.Info, "Creating server %s", p.Id())
	p.RunningEnvironment.DisplayToConsole(true, "Allocating server\n")
	err = p.RunningEnvironment.Create()
	if err != nil {
		p.Log(logging.Error, "Error creating server: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Failed to create server\n")
	} else {
		p.RunningEnvironment.DisplayToConsole(true, "Server allocated\n")
	}

	return
}

// Destroys the server.
// This will delete the server, environment, and any files related to it.
func (p *Program) Destroy() (err error) {
	p.Log(logging.Info, "Destroying server %s", p.Id())
	process, err := operations.GenerateProcess(p.Uninstallation, p.RunningEnvironment, p.DataToMap(), p.Execution.EnvironmentVariables)
	if err != nil {
		p.Log(logging.Error, "Error uninstalling server: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Failed to uninstall server\n")
		return
	}

	err = process.Run(p.RunningEnvironment)
	if err != nil {
		p.Log(logging.Error, "Error uninstalling server: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Failed to uninstall server\n")
		return
	}

	err = p.RunningEnvironment.Delete()
	if err != nil {
		p.Log(logging.Error, "Error uninstalling server: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Failed to uninstall server\n")
	}
	err = p.Scheduler.Rebuild()
	if err != nil {
		p.Log(logging.Error, "Error uninstalling server: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Failed to uninstall server\n%s\n", err.Error())
	}
	err = p.Scheduler.Rebuild()
	if err != nil {
		p.Log(logging.Error, "Error uninstalling server: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Failed to uninstall server\n%s\n", err.Error())
	}
	return
}

func (p *Program) Install() (err error) {
	if !p.IsEnabled() {
		p.Log(logging.Error, "Server %s is not enabled, cannot install", p.Id())
		return pufferpanel.ErrServerDisabled
	}

	p.Log(logging.Info, "Installing server %s", p.Id())
	running, err := p.IsRunning()
	if err != nil {
		p.Log(logging.Error, "Error checking server status: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Error on checking to see if server is running\n")
		return
	}

	if running {
		err = p.Stop()
	}

	if err != nil {
		p.Log(logging.Error, "Error stopping server: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Failed to stop server\n")
		return
	}

	p.RunningEnvironment.DisplayToConsole(true, "Installing server\n")

	err = os.MkdirAll(p.RunningEnvironment.GetRootDirectory(), 0755)
	if err != nil && !os.IsExist(err) {
		p.Log(logging.Error, "Error creating server directory: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Failed to create server directory\n")
		return
	}

	if len(p.Installation) > 0 {
		process, err := operations.GenerateProcess(p.Installation, p.GetEnvironment(), p.DataToMap(), p.Execution.EnvironmentVariables)
		if err != nil {
			p.Log(logging.Error, "Error installing server: %s", err)
			p.RunningEnvironment.DisplayToConsole(true, "Failed to install server\n")
			return err
		}

		err = process.Run(p.RunningEnvironment)
		if err != nil {
			p.Log(logging.Error, "Error installing server: %s", err)
			p.RunningEnvironment.DisplayToConsole(true, "Failed to install server\n")
			return err
		}
	}

	p.RunningEnvironment.DisplayToConsole(true, "Server installed\n")
	return
}

func (p *Program) IsRunning() (isRunning bool, err error) {
	isRunning, err = p.RunningEnvironment.IsRunning()
	return
}

func (p *Program) Execute(command string) (err error) {
	err = p.RunningEnvironment.ExecuteInMainProcess(command)
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
	p.RunningEnvironment = environment
	return
}

func (p *Program) Id() string {
	return p.Identifier
}

func (p *Program) GetEnvironment() pufferpanel.Environment {
	return p.RunningEnvironment
}

func (p *Program) SetAutoStart(isAutoStart bool) (err error) {
	p.Execution.AutoStart = isAutoStart
	return
}

func (p *Program) IsAutoStart() (isAutoStart bool) {
	isAutoStart = p.Execution.AutoStart
	return
}

func (p *Program) Save() (err error) {
	p.Log(logging.Info, "Saving server %s", p.Id())

	file := filepath.Join(config.ServersFolder.Value(), p.Id()+".json")

	if !p.valid() {
		p.Log(logging.Error, "Server %s contained invalid data, this server is.... broken", p.Identifier)
		//we can't even reload from disk....
		//so, puke back, and for now we'll handle it later
		return pufferpanel.ErrUnknownError
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(file, data, 0664)
	return
}

func (p *Program) EditData(data map[string]pufferpanel.Variable, overrideUser bool) (err error) {
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
		ip = cast.ToString(ipData.Value)
	}

	if portData, ok := data["port"]; ok {
		port = cast.ToString(portData.Value)
	}

	return ip + ":" + port
}

func (p *Program) afterExit(graceful bool) {
	if graceful {
		p.CrashCounter = 0
	}

	mapping := p.DataToMap()
	mapping["success"] = graceful

	processes, err := operations.GenerateProcess(p.Execution.PostExecution, p.RunningEnvironment, mapping, p.Execution.EnvironmentVariables)
	if err != nil {
		p.Log(logging.Error, "Error running post processing for server %s: %s", p.Id(), err)
		p.RunningEnvironment.DisplayToConsole(true, "Failed to run post-execution steps\n")
		return
	}
	p.RunningEnvironment.DisplayToConsole(true, "Running post-execution steps\n")
	p.Log(logging.Info, "Running post execution steps: %s", p.Id())

	err = processes.Run(p.RunningEnvironment)
	if err != nil {
		p.Log(logging.Error, "Error running post processing for server: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Failed to run post-execution steps\n")
		return
	}

	if !p.Execution.AutoRestartFromCrash && !p.Execution.AutoRestartFromGraceful {
		return
	}

	if graceful && p.Execution.AutoRestartFromGraceful {
		StartViaService(p)
	} else if !graceful && p.Execution.AutoRestartFromCrash && p.CrashCounter < config.CrashLimit.Value() {
		p.CrashCounter++
		StartViaService(p)
	}
}

func (p *Program) GetItem(name string) (*FileData, error) {
	info, err := p.GetFileServer().Stat(name)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		files, _ := p.GetFileServer().ReadDir(name)
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

		for i, file := range files {
			newFile := messages.FileDesc{
				Name: file.Name(),
				File: !file.IsDir(),
			}

			if !file.IsDir() && file.Type()&os.ModeSymlink == 0 {
				infoData, err := p.GetFileServer().Stat(filepath.Join(name, file.Name()))
				if err != nil {
					continue
				}
				newFile.Size = infoData.Size()
				newFile.Modified = infoData.ModTime().Unix()
				newFile.Extension = filepath.Ext(file.Name())
			}

			fileNames[i+offset] = newFile
		}

		return &FileData{FileList: fileNames}, nil
	} else {
		file, err := p.GetFileServer().Open(name)
		if err != nil {
			return nil, err
		}
		return &FileData{Contents: file, ContentLength: info.Size(), Name: info.Name()}, nil
	}
}

func (p *Program) ArchiveItems(files []string, destination string) error {
	var targets []string
	for _, name := range files {
		targetFile := pufferpanel.JoinPath(p.GetEnvironment().GetRootDirectory(), name)
		if !pufferpanel.EnsureAccess(targetFile, p.GetEnvironment().GetRootDirectory()) {
			return pufferpanel.ErrIllegalFileAccess
		}
		targets = append(targets, targetFile)
	}

	destination = pufferpanel.JoinPath(p.GetEnvironment().GetRootDirectory(), destination)
	if !pufferpanel.EnsureAccess(destination, p.GetEnvironment().GetRootDirectory()) {
		return pufferpanel.ErrIllegalFileAccess
	}

	// This may technically error out in other cases
	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		return pufferpanel.ErrFileExists
	}
	return archiver.Archive(targets, destination)
}

func (p *Program) Extract(source, destination string) error {
	sourceFile := pufferpanel.JoinPath(p.GetEnvironment().GetRootDirectory(), source)
	destinationFile := pufferpanel.JoinPath(p.GetEnvironment().GetRootDirectory(), destination)

	if !pufferpanel.EnsureAccess(sourceFile, p.GetEnvironment().GetRootDirectory()) || !pufferpanel.EnsureAccess(destinationFile, p.GetEnvironment().GetRootDirectory()) {
		return pufferpanel.ErrIllegalFileAccess
	}

	// destination shouldn't exist
	if _, err := os.Stat(destinationFile); os.IsNotExist(err) {
		return pufferpanel.ErrFileExists
	}

	return archiver.Unarchive(sourceFile, destinationFile)
}

func (p *Program) ExecuteTask(task pufferpanel.Task) (err error) {
	ops := task.Operations
	if len(ops) > 0 {
		p.RunningEnvironment.DisplayToConsole(true, "Running task %s\n", task.Name)
		var process operations.OperationProcess
		process, err = operations.GenerateProcess(ops, p.GetEnvironment(), p.DataToMap(), p.Execution.EnvironmentVariables)
		if err != nil {
			p.Log(logging.Error, "Error setting up tasks: %s", err)
			p.RunningEnvironment.DisplayToConsole(true, "Failed to setup tasks\n")
			p.RunningEnvironment.DisplayToConsole(true, "%s\n", err.Error())
			return
		}

		err = process.Run(p.RunningEnvironment)
		if err != nil {
			p.Log(logging.Error, "Error setting up tasks: %s", err)
			p.RunningEnvironment.DisplayToConsole(true, "Failed to setup tasks\n")
			p.RunningEnvironment.DisplayToConsole(true, "%s\n", err.Error())
			return
		}
		p.RunningEnvironment.DisplayToConsole(true, "Task %s finished\n", task.Name)
	}
	return
}

func (p *Program) GetFileServer() pufferpanel.FileServer {
	return p.fs
}

func (p *Program) valid() bool {
	//we need a type at least, this is a safe check
	if p.Type.Type == "" {
		return false
	}

	//check the env object, if it's even checkable
	env, casted := p.Environment.(map[string]interface{})
	if !casted || len(env) == 0 {
		return false
	}

	v, exists := env["type"]
	if !exists {
		return false
	}

	str, ok := v.(string)
	if !ok || str == "" {
		return false
	}

	return true
}

func (p *Program) Log(l *log.Logger, format string, obj ...interface{}) {
	msg := fmt.Sprintf("[%s] ", p.Id()) + format
	l.Printf(msg, obj...)
}
