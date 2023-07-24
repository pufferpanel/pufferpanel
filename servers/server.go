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

package servers

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/mholt/archiver/v3"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/conditions"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/messages"
	"github.com/spf13/cast"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Server struct {
	pufferpanel.Server

	CrashCounter       int                     `json:"-"`
	RunningEnvironment pufferpanel.Environment `json:"-"`
	Scheduler          Scheduler               `json:"-"`
	stopChan           chan bool
	waitForConsole     sync.Locker
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

func StartViaService(p *Server) {
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
		program := next.Value.(*Server)
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

func (p *Server) DataToMap() map[string]interface{} {
	var result = p.Server.DataToMap()
	result["rootDir"] = p.RunningEnvironment.GetRootDirectory()
	result["core:os"] = runtime.GOOS
	result["core:arch"] = runtime.GOARCH

	return result
}

func CreateProgram() *Server {
	p := &Server{
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
			Groups:         make([]pufferpanel.Group, 0),
		},
	}
	p.Scheduler = NewScheduler(p)
	p.stopChan = make(chan bool)
	p.waitForConsole = &sync.Mutex{}
	return p
}

// Starts the program.
// This includes starting the environment if it is not running.
func (p *Server) Start() error {
	if !p.IsEnabled() {
		p.Log(logging.Error, "Server %s is not enabled, cannot start", p.Id())
		return pufferpanel.ErrServerDisabled
	}
	if r, err := p.IsRunning(); r || err != nil {
		if err == nil {
			err = pufferpanel.ErrServerRunning
		}
		return err
	}

	p.Log(logging.Info, "Starting server %s", p.Id())
	p.RunningEnvironment.DisplayToConsole(true, "Starting server\n")

	data := p.DataToMap()

	process, err := GenerateProcess(p.Execution.PreExecution, p.RunningEnvironment, data, p.Execution.EnvironmentVariables)
	if err != nil {
		p.Log(logging.Error, "Error generating pre-execution steps: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Error running pre execute\n")
		return err
	}

	err = process.Run(p)
	if err != nil {
		p.Log(logging.Error, "Error running pre-execution steps: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Error running pre execute\n")
		return err
	}

	var command string

	if c, ok := p.Execution.Command.(string); ok {
		command = c
	} else {
		//we have a list
		var possibleCommands []pufferpanel.Command
		err = pufferpanel.UnmarshalTo(p.Execution.Command, &possibleCommands)
		if err != nil {
			return err
		}

		var defaultCommand pufferpanel.Command
		var commandToRun pufferpanel.Command
		for _, v := range possibleCommands {
			if v.If == "" {
				defaultCommand = v
				break
			}
		}

		for _, v := range possibleCommands {
			//now... we see which command to use
			if v.If == "" {
				continue
			}
			useThis, err := p.RunCondition(v.If, nil)
			if err != nil {
				p.Log(logging.Error, "error starting server %s: %s", p.Id(), err)
				p.RunningEnvironment.DisplayToConsole(true, " Failed to start server\n")
				return err
			}
			if useThis {
				commandToRun = v
				break
			}
		}

		command = commandToRun.Command

		//if no command, use default
		if command == "" {
			command = defaultCommand.Command
		}
	}

	if command == "" {
		err = pufferpanel.ErrNoCommand
		p.Log(logging.Error, "error starting server %s: %s", p.Id(), err)
		p.RunningEnvironment.DisplayToConsole(true, " Failed to start server\n")
		return err
	}

	commandLine := pufferpanel.ReplaceTokens(command, data)
	if p.Execution.WorkingDirectory == "${rootDir}" {
		p.Execution.WorkingDirectory = ""
	}
	workDir := pufferpanel.ReplaceTokens(p.Execution.WorkingDirectory, data)

	if !pufferpanel.EnsureAccess(path.Join(p.RunningEnvironment.GetRootDirectory(), workDir), p.RunningEnvironment.GetRootDirectory()) {
		p.Log(logging.Error, "Working directory is invalid for server: %s", workDir)
		p.RunningEnvironment.DisplayToConsole(true, "Working directory is invalid for server: %s", workDir)
		return err
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
		return err
	}

	//server started, now kick off our special "hook" for the console
	consoleConfig := p.GetEnvironment().GetStdOutConfiguration()
	switch consoleConfig.Type {
	case "file":
		{
			go func() {
				defer p.waitForConsole.Unlock()
				p.waitForConsole.Lock()
				//now... we need to wait for a file to exist
				//once it's there, we read until we die
				var file *os.File
				for file, err = os.Open(consoleConfig.File); os.IsNotExist(err); {
				}
				_, _ = io.Copy(p.RunningEnvironment.GetWrapper(), file)
			}()
		}
	}

	return err
}

// Stops the program.
// This will also stop the environment it is ran in.
func (p *Server) Stop() error {
	var err error
	if r, err := p.IsRunning(); !r || err != nil {
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
	return err
}

// Kills the program.
// This will also stop the environment it is ran in.
func (p *Server) Kill() (err error) {
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
func (p *Server) Create() (err error) {
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
func (p *Server) Destroy() (err error) {
	p.Log(logging.Info, "Destroying server %s", p.Id())
	process, err := GenerateProcess(p.Uninstallation, p.RunningEnvironment, p.DataToMap(), p.Execution.EnvironmentVariables)
	if err != nil {
		p.Log(logging.Error, "Error uninstalling server: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Failed to uninstall server\n")
		return
	}

	err = process.Run(p)
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

func (p *Server) Install() error {
	if !p.IsEnabled() {
		p.Log(logging.Error, "Server %s is not enabled, cannot install", p.Id())
		return pufferpanel.ErrServerDisabled
	}

	p.Log(logging.Info, "Installing server %s", p.Id())
	r, err := p.IsRunning()
	if err != nil {
		p.Log(logging.Error, "Error checking server status: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Error on checking to see if server is running\n")
		return err
	}

	if r {
		err = p.Stop()
	}

	if err != nil {
		p.Log(logging.Error, "Error stopping server: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Failed to stop server\n")
		return err
	}

	p.RunningEnvironment.DisplayToConsole(true, "Installing server\n")

	err = os.MkdirAll(p.RunningEnvironment.GetRootDirectory(), 0755)
	if err != nil && !os.IsExist(err) {
		p.Log(logging.Error, "Error creating server directory: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Failed to create server directory\n")
		return err
	}

	if len(p.Installation) > 0 {
		var process OperationProcess
		process, err = GenerateProcess(p.Installation, p.GetEnvironment(), p.DataToMap(), p.Execution.EnvironmentVariables)
		if err != nil {
			p.Log(logging.Error, "Error installing server: %s", err)
			p.RunningEnvironment.DisplayToConsole(true, "Failed to install server\n")
			return err
		}

		err = process.Run(p)
		if err != nil {
			p.Log(logging.Error, "Error installing server: %s", err)
			p.RunningEnvironment.DisplayToConsole(true, "Failed to install server\n")
			return err
		}
	}

	p.RunningEnvironment.DisplayToConsole(true, "Server installed\n")
	return nil
}

func (p *Server) IsRunning() (bool, error) {
	return p.RunningEnvironment.IsRunning()
}

func (p *Server) Execute(command string) (err error) {
	err = p.RunningEnvironment.ExecuteInMainProcess(command)
	return
}

func (p *Server) SetEnabled(isEnabled bool) (err error) {
	p.Execution.Disabled = !isEnabled
	return
}

func (p *Server) IsEnabled() (isEnabled bool) {
	return !p.Execution.Disabled
}

func (p *Server) SetEnvironment(environment pufferpanel.Environment) (err error) {
	p.RunningEnvironment = environment
	return
}

func (p *Server) Id() string {
	return p.Identifier
}

func (p *Server) GetEnvironment() pufferpanel.Environment {
	return p.RunningEnvironment
}

func (p *Server) SetAutoStart(isAutoStart bool) (err error) {
	p.Execution.AutoStart = isAutoStart
	return
}

func (p *Server) IsAutoStart() (isAutoStart bool) {
	isAutoStart = p.Execution.AutoStart
	return
}

func (p *Server) Save() (err error) {
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

	err = os.WriteFile(file, data, 0664)
	return
}

func (p *Server) EditData(data map[string]pufferpanel.Variable, overrideUser bool) (err error) {
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

	err = p.Save()
	return
}

func (p *Server) GetData() map[string]pufferpanel.Variable {
	return p.Variables
}

func (p *Server) GetNetwork() string {
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

func (p *Server) afterExit(exitCode int) {
	graceful := exitCode == 0
	if graceful {
		p.CrashCounter = 0
	}

	mapping := p.DataToMap()
	mapping["success"] = graceful
	mapping["exitCode"] = exitCode

	processes, err := GenerateProcess(p.Execution.PostExecution, p.RunningEnvironment, mapping, p.Execution.EnvironmentVariables)
	if err != nil {
		p.Log(logging.Error, "Error running post processing for server %s: %s", p.Id(), err)
		p.RunningEnvironment.DisplayToConsole(true, "Failed to run post-execution steps\n")
		return
	}
	p.RunningEnvironment.DisplayToConsole(true, "Running post-execution steps\n")
	p.Log(logging.Info, "Running post execution steps: %s", p.Id())

	err = processes.Run(p)
	if err != nil {
		p.Log(logging.Error, "Error running post processing for server: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Failed to run post-execution steps\n")
		return
	}

	if !p.Execution.AutoRestartFromCrash && !p.Execution.AutoRestartFromGraceful {
		return
	}

	p.stopChan <- true
	//wait for the console to be done with it's work, if it has any
	p.waitForConsole.Lock()
	p.waitForConsole.Unlock()

	if graceful && p.Execution.AutoRestartFromGraceful {
		StartViaService(p)
	} else if !graceful && p.Execution.AutoRestartFromCrash && p.CrashCounter < config.CrashLimit.Value() {
		p.CrashCounter++
		StartViaService(p)
	}
}

func (p *Server) GetItem(name string) (*FileData, error) {
	targetFile := pufferpanel.JoinPath(p.GetEnvironment().GetRootDirectory(), name)
	if !pufferpanel.EnsureAccess(targetFile, p.GetEnvironment().GetRootDirectory()) {
		return nil, pufferpanel.ErrIllegalFileAccess
	}

	info, err := os.Stat(targetFile)

	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		files, _ := os.ReadDir(targetFile)
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
				infoData, _ := file.Info()
				newFile.Size = infoData.Size()
				newFile.Modified = infoData.ModTime().Unix()
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

func (p *Server) CreateFolder(name string) error {
	folder := pufferpanel.JoinPath(p.GetEnvironment().GetRootDirectory(), name)

	if !pufferpanel.EnsureAccess(folder, p.GetEnvironment().GetRootDirectory()) {
		return pufferpanel.ErrIllegalFileAccess
	}
	return os.MkdirAll(folder, 0755)
}

func (p *Server) OpenFile(name string) (io.WriteCloser, error) {
	targetFile := pufferpanel.JoinPath(p.GetEnvironment().GetRootDirectory(), name)

	if !pufferpanel.EnsureAccess(targetFile, p.GetEnvironment().GetRootDirectory()) {
		return nil, pufferpanel.ErrIllegalFileAccess
	}

	file, err := os.Create(targetFile)
	return file, err
}

func (p *Server) DeleteItem(name string) error {
	targetFile := pufferpanel.JoinPath(p.GetEnvironment().GetRootDirectory(), name)

	if !pufferpanel.EnsureAccess(targetFile, p.GetEnvironment().GetRootDirectory()) {
		return pufferpanel.ErrIllegalFileAccess
	}

	return os.RemoveAll(targetFile)
}

func (p *Server) ArchiveItems(files []string, destination string) error {
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

func (p *Server) Extract(source, destination string) error {
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

func (p *Server) ExecuteTask(task pufferpanel.Task) (err error) {
	ops := task.Operations
	if len(ops) > 0 {
		p.RunningEnvironment.DisplayToConsole(true, "Running task %s\n", task.Name)
		var process OperationProcess
		process, err = GenerateProcess(ops, p.GetEnvironment(), p.DataToMap(), p.Execution.EnvironmentVariables)
		if err != nil {
			p.Log(logging.Error, "Error setting up tasks: %s", err)
			p.RunningEnvironment.DisplayToConsole(true, "Failed to setup tasks\n")
			p.RunningEnvironment.DisplayToConsole(true, "%s\n", err.Error())
			return
		}

		err = process.Run(p)
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

func (p *Server) valid() bool {
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

func (p *Server) Log(l *log.Logger, format string, obj ...interface{}) {
	msg := fmt.Sprintf("[%s] ", p.Id()) + format
	l.Printf(msg, obj...)
}

func (p *Server) RunCondition(condition interface{}, extraData map[string]interface{}) (bool, error) {
	data := map[string]interface{}{
		conditions.VariableEnv:      p.RunningEnvironment.GetBase().Type,
		conditions.VariableServerId: p.Id(),
	}

	if extraData != nil {
		for k, v := range extraData {
			data[k] = v
		}
	}

	if p.Variables != nil {
		for k, v := range p.Variables {
			data[k] = v.Value
		}
	}

	return conditions.ResolveIf(condition, data, CreateFunctions(p.GetEnvironment()))
}
