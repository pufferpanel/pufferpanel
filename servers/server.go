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
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Server struct {
	pufferpanel.Server

	CrashCounter       int                     `json:"-"`
	RunningEnvironment pufferpanel.Environment `json:"-"`
	Scheduler          *Scheduler              `json:"-"`
	stopChan           chan bool
	waitForConsole     sync.Locker
	fileServer         pufferpanel.FileServer
}

var queue *list.List
var lock = sync.Mutex{}
var startQueueTicker, statTicker *time.Ticker
var running = false

func init() {
	archiver.DefaultZip.OverwriteExisting = true
	archiver.DefaultTarGz.OverwriteExisting = true
}

func InitService() {
	queue = list.New()
	running = true
	go processQueue()
	go processStats()
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
	startQueueTicker.Stop()
	statTicker.Stop()
}

func processQueue() {
	startQueueTicker = time.NewTicker(time.Second)
	for range startQueueTicker.C {
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

func processStats() {
	statTicker = time.NewTicker(5 * time.Second)
	for range statTicker.C {
		SendStatsForServers()
	}
}

func SendStatsForServers() {
	var wg sync.WaitGroup
	for _, v := range allServers {
		wg.Add(1)
		go func(p *Server) {
			defer wg.Done()
			stats, err := p.GetEnvironment().GetStats()
			if err != nil {
				return
			}
			_ = p.GetEnvironment().GetStatsTracker().WriteMessage(&messages.Stat{
				Memory: stats.Memory,
				Cpu:    stats.Cpu,
			})
		}(v)
	}
	wg.Wait()
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
				AutoStart:               false,
				AutoRestartFromCrash:    false,
				AutoRestartFromGraceful: false,
				PreExecution:            make([]pufferpanel.ConditionalMetadataType, 0),
				PostExecution:           make([]pufferpanel.ConditionalMetadataType, 0),
				EnvironmentVariables:    make(map[string]string),
			},
			Type:           pufferpanel.Type{Type: "standard"},
			Variables:      make(map[string]pufferpanel.Variable),
			Display:        "Unknown server",
			Installation:   make([]pufferpanel.ConditionalMetadataType, 0),
			Uninstallation: make([]pufferpanel.ConditionalMetadataType, 0),
			Groups:         make([]pufferpanel.Group, 0),
		},
	}
	p.stopChan = make(chan bool)
	p.waitForConsole = &sync.Mutex{}
	return p
}

// Start Starts the program.
// This includes starting the environment if it is not running.
func (p *Server) Start() error {
	if r, err := p.IsRunning(); r || err != nil {
		if err == nil {
			err = pufferpanel.ErrServerRunning
		}
		return err
	}

	p.Log(logging.Info, "Starting server %s", p.Id())
	p.RunningEnvironment.DisplayToConsole(true, "Starting server\n")

	process, err := GenerateProcess(p.Execution.PreExecution, p.RunningEnvironment, p.DataToMap(), p.Execution.EnvironmentVariables)
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

	var command pufferpanel.Command

	if c, ok := p.Execution.Command.(string); ok {
		command = pufferpanel.Command{Command: c}
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

		command = commandToRun

		//if no command, use default
		if command.Command == "" {
			command = defaultCommand
		}
	}

	if command.StdIn.Type == "" {
		command.StdIn = p.Execution.Stdin
	}

	data := p.DataToMap()

	commandLine := pufferpanel.ReplaceTokens(command.Command, data)

	cmd, args := pufferpanel.SplitArguments(commandLine)
	err = p.RunningEnvironment.ExecuteAsync(pufferpanel.ExecutionData{
		Command:     cmd,
		Arguments:   args,
		Environment: pufferpanel.ReplaceTokensInMap(p.Execution.EnvironmentVariables, data),
		Variables:   p.DataToMap(),
		Callback:    p.afterExit,
		StdInConfig: command.StdIn,
	})

	if err != nil {
		p.Log(logging.Error, "error starting server %s: %s", p.Id(), err)
		p.RunningEnvironment.DisplayToConsole(true, " Failed to start server\n")
		return err
	}

	return err
}

// Stop Stops the program.
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

// Kill Kills the program.
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

// Create Creates any files needed for the program.
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

// Destroy Destroys the server.
// This will delete the server, environment, and any files related to it.
func (p *Server) Destroy() (err error) {
	p.Log(logging.Info, "Destroying server %s", p.Id())

	if p.Scheduler != nil {
		p.Scheduler.Stop()
	}

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

	return
}

func (p *Server) Install() error {
	if p.GetEnvironment().IsInstalling() {
		return nil
	}

	p.GetEnvironment().SetInstalling(true)
	defer p.GetEnvironment().SetInstalling(false)

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

		data := p.DataToMap()
		process, err = GenerateProcess(p.Installation, p.RunningEnvironment, data, p.Execution.EnvironmentVariables)
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

func (p *Server) EditData(data map[string]interface{}) (err error) {
	for k, v := range data {
		var elem pufferpanel.Variable

		if _, ok := p.Variables[k]; ok {
			elem = p.Variables[k]
		}
		if !elem.UserEditable {
			continue
		}

		elem.Value = v

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
	graceful := exitCode == p.Execution.ExpectedExitCode
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
	/*if _, err := os.Stat(destinationFile); os.IsNotExist(err) {
		return pufferpanel.ErrFileExists
	}*/
	return archiver.Unarchive(sourceFile, destinationFile)
}

func (p *Server) valid() bool {
	//we need a type at least, this is a safe check
	if p.Type.Type == "" {
		return false
	}

	if p.Environment.Type == "" {
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

	for k, v := range extraData {
		data[k] = v
	}

	if p.Variables != nil {
		for k, v := range p.Variables {
			data[k] = v.Value
		}
	}

	return conditions.ResolveIf(condition, data, CreateFunctions(p.GetEnvironment()))
}

func (p *Server) GetFileServer() pufferpanel.FileServer {
	return p.fileServer
}
