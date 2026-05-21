package servers

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/mholt/archiver/v3"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/conditions"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/files"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/utils"
	"github.com/spf13/cast"
)

type Server struct {
	pufferpanel.DaemonServer `json:"-"`
	pufferpanel.Server

	CrashCounter       int                      `json:"-"`
	RunningEnvironment *pufferpanel.Environment `json:"-"`
	Scheduler          *Scheduler               `json:"-"`
	stopChan           chan bool
	waitForConsole     sync.Locker
	fileServer         files.FileServer
	backupFileServer   files.FileServer
	backingUp          bool
	restoring          bool
	keepAlive          *time.Ticker
	keepAliveChan      chan bool
	isUnsafeRunning    *sync.Mutex
}

var queue *list.List
var lock = sync.Mutex{}
var startQueueTicker, statTicker *time.Ticker
var running = false

var ErrServerTypeRequired = errors.New("server type is required")
var ErrEnvironmentTypeRequired = errors.New("environment type is required")

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

			_ = p.GetEnvironment().GetStatsTracker().WriteMessage(pufferpanel.Transmission{
				Message: stats,
				Type:    pufferpanel.MessageTypeStats,
			})
		}(v)
	}
	wg.Wait()
}

type FileData struct {
	Contents      io.ReadCloser
	ContentLength int64
	FileList      []pufferpanel.FileDesc
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
	p.stopChan = make(chan bool, 1)
	p.waitForConsole = &sync.Mutex{}
	p.isUnsafeRunning = &sync.Mutex{}
	p.keepAliveChan = make(chan bool)
	return p
}

// Start Starts the program.
// This includes starting the environment if it is not running.
func (p *Server) Start() (err error) {
	if err = p.IsIdle(); err != nil {
		return
	}

	if !p.isUnsafeRunning.TryLock() {
		return pufferpanel.ErrServerRunning
	}

	defer func() {
		if err != nil {
			//since we had an issue starting, we need to "unlock"
			//otherwise, the afterExit will unlock this at the end
			p.isUnsafeRunning.Unlock()
		}
	}()

	p.Log(logging.Info, "Starting server %s", p.Id())
	p.RunningEnvironment.DisplayToConsole(true, "Starting server\n")

	var process OperationProcess
	process, err = GenerateProcess(p.Execution.PreExecution, p.RunningEnvironment, p.DataToMap(), p.Execution.EnvironmentVariables)
	if err != nil {
		p.Log(logging.Error, "Error generating pre-execution steps: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Error running pre execute\n")
		return
	}

	err = process.Run(p)
	if err != nil {
		p.Log(logging.Error, "Error running pre-execution steps: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Error running pre execute\n")
		return
	}

	var command pufferpanel.Command

	if c, ok := p.Execution.Command.(string); ok {
		command = pufferpanel.Command{Command: c}
	} else {
		//we have a list
		var possibleCommands []pufferpanel.Command
		err = utils.UnmarshalTo(p.Execution.Command, &possibleCommands)
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

	commandLine := utils.ReplaceTokens(command.Command, data)
	err = p.RunningEnvironment.ExecuteAsync(pufferpanel.ExecutionData{
		WorkingDirectory: p.Server.Execution.WorkingDirectory,
		Command:          commandLine,
		Environment:      utils.ReplaceTokensInMap(p.Execution.EnvironmentVariables, data),
		Variables:        p.DataToMap(),
		Callback:         p.afterExit,
		StdInConfig:      command.StdIn,
	})

	if err != nil {
		p.Log(logging.Error, "error starting server %s: %s", p.Id(), err)
		p.RunningEnvironment.DisplayToConsole(true, " Failed to start server\n")
		return err
	}

	//keepalive!
	if p.KeepAlive.Frequency != "" && p.KeepAlive.Command != "" {
		dur, err := time.ParseDuration(p.KeepAlive.Frequency)
		if err != nil {
			p.RunningEnvironment.DisplayToConsole(true, " Failed to enable keep-alive: %s", err)
			return nil
		}
		if p.keepAlive == nil {
			p.keepAlive = time.NewTicker(dur)
		} else {
			p.keepAlive.Reset(dur)
		}

		if p.keepAlive != nil {
			go func() {
				for {
					select {
					case <-p.keepAliveChan:
						return
					case <-p.keepAlive.C:
						_ = p.RunningEnvironment.ExecuteInMainProcess(p.KeepAlive.Command)
					}
				}
			}()
		}
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
	if err := p.IsIdle(); err != nil {
		return err
	}

	p.Log(logging.Info, "Destroying server %s", p.Id())

	p.Log(logging.Debug, "Stopping scheduler")
	if p.Scheduler != nil && p.Scheduler.IsRunning() {
		p.Scheduler.Stop()
	}

	p.Log(logging.Debug, "Starting uninstall processes")
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

	p.Log(logging.Debug, "Deleting environment")
	err = p.RunningEnvironment.Delete()
	if err != nil {
		p.Log(logging.Error, "Error uninstalling server: %s", err)
		p.RunningEnvironment.DisplayToConsole(true, "Failed to uninstall server\n")
	}

	return
}

func (p *Server) Install() error {
	if err := p.IsIdle(); err != nil {
		return err
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

func (p *Server) SetEnvironment(environment *pufferpanel.Environment) (err error) {
	p.RunningEnvironment = environment
	return
}

func (p *Server) Id() string {
	return p.Identifier
}

func (p *Server) GetEnvironment() *pufferpanel.Environment {
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

	if err = p.valid(); err != nil {
		p.Log(logging.Error, "Server %s contained invalid data, this server is.... broken", p.Identifier)
		//we can't even reload from disk....
		//so, puke back, and for now we'll handle it later
		return err
	}

	var data []byte
	data, err = json.MarshalIndent(p, "", "  ")
	if err != nil {
		return
	}

	err = p.GetFileServer().WriteFile(file, data)
	return
}

func (p *Server) EditData(data map[string]interface{}, asAdmin bool) (err error) {
	for k, v := range data {
		var elem pufferpanel.Variable

		if _, ok := p.Variables[k]; ok {
			elem = p.Variables[k]
		}
		if !asAdmin && !elem.UserEditable {
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
	if p.keepAlive != nil {
		p.keepAlive.Stop()
		p.keepAliveChan <- true
	}

	graceful := exitCode == p.Execution.ExpectedExitCode
	if graceful {
		p.CrashCounter = 0
	}

	mapping := p.DataToMap()
	mapping["success"] = graceful
	mapping["exitCode"] = exitCode

	processes, err := GenerateProcess(p.Execution.PostExecution, p.RunningEnvironment, mapping, p.Execution.EnvironmentVariables)
	if err != nil {
		p.Log(logging.Error, "Error generating post processing tasks for server %s: %s", p.Id(), err)
		p.RunningEnvironment.DisplayToConsole(true, "Failed to run post-execution steps\n")
	} else {
		p.RunningEnvironment.DisplayToConsole(true, "Running post-execution steps\n")
		p.Log(logging.Info, "Running post execution steps: %s", p.Id())

		err = processes.Run(p)
		if err != nil {
			p.Log(logging.Error, "Error running post processing for server: %s", err)
			p.RunningEnvironment.DisplayToConsole(true, "Failed to run post-execution steps\n")
		}
	}

	p.isUnsafeRunning.Unlock()

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
		fileList, _ := p.GetFileServer().ReadDir(name)
		var fileNames []pufferpanel.FileDesc
		offset := 0
		if name == "" || name == "." || name == "/" {
			fileNames = make([]pufferpanel.FileDesc, len(fileList))
		} else {
			fileNames = make([]pufferpanel.FileDesc, len(fileList)+1)
			fileNames[0] = pufferpanel.FileDesc{
				Name: "..",
				File: false,
			}
			offset = 1
		}

		//validate any symlinks are valid

		for i, file := range fileList {
			newFile := pufferpanel.FileDesc{
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

func (p *Server) ArchiveItems(sourceFiles []string, destination string) error {
	// This may technically error out in other cases
	if _, err := p.GetFileServer().Stat(destination); err != nil && !os.IsNotExist(err) {
		return pufferpanel.ErrFileExists
	}
	return files.Compress(p.GetFileServer(), destination, p.GetFileServer(), sourceFiles)
}

func (p *Server) Extract(source, destination string) error {
	return files.Extract(p.GetFileServer(), source, p.GetFileServer(), destination, "*")
}

func (p *Server) StartBackup() (string, error) {
	if err := p.IsIdle(); err != nil {
		return "", err
	}

	p.backingUp = true
	c := make(chan bool)
	go func(d chan bool) {
		r := <-d
		p.backingUp = false
		if r {
			p.RunningEnvironment.DisplayToConsole(true, "Backup complete")
		} else {
			p.RunningEnvironment.DisplayToConsole(true, "Backup failed")
		}
	}(c)

	p.RunningEnvironment.DisplayToConsole(true, "Backing up server")

	backupId, err := uuid.NewV4()
	if err != nil {
		c <- false
		return "", err
	}
	backupFileName := backupId.String() + ".tar.gz"

	go func(file string, d chan bool) {
		defer func() {
			d <- true
		}()
		sourceFiles := []string{filepath.Join(p.GetFileServer().Prefix())}

		err = files.Compress(p.GetBackupFileServer(), backupFileName, p.fileServer, sourceFiles)
		if err != nil {
			p.Log(logging.Error, "Error creating backup file: %s", err)
			p.RunningEnvironment.DisplayToConsole(true, "Failed to create backup file")
		}
	}(backupFileName, c)

	return backupFileName, nil
}

func (p *Server) DeleteBackup(fileName string) error {
	err := p.GetBackupFileServer().Remove(fileName)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func (p *Server) StartRestore(fileName string) error {
	if err := p.IsIdle(); err != nil {
		return err
	}

	p.restoring = true
	c := make(chan bool)
	go func(d chan bool) {
		r := <-d
		p.restoring = false
		if r {
			p.RunningEnvironment.DisplayToConsole(true, "Restore complete")
		} else {
			p.RunningEnvironment.DisplayToConsole(true, "Restore failed")
		}
	}(c)

	p.RunningEnvironment.DisplayToConsole(true, "Restoring server")

	_, err := p.GetBackupFileServer().Stat(fileName)
	if err != nil && !os.IsNotExist(err) {
		c <- false
		return err
	}

	go func(source string, d chan bool) {
		defer func() {
			d <- true
		}()

		//Check if any files exist, as remove all errors if its empty
		existingFiles, err := p.GetFileServer().Glob("*")
		if err != nil {
			p.Log(logging.Error, "Error globbing files: %s", err)
			return
		}

		for _, existingFile := range existingFiles {
			file, err := p.GetFileServer().Stat(existingFile)
			if err != nil {
				p.Log(logging.Error, "Error deleting files: %s", err)
				return
			}

			if file.IsDir() {
				err = p.GetFileServer().RemoveAll(existingFile)
			} else {
				err = p.GetFileServer().Remove(existingFile)
			}

			if err != nil {
				p.Log(logging.Error, "Error deleting files: %s", err)
				return
			}
		}

		err = files.Extract(p.GetBackupFileServer(), source, p.GetFileServer(), "", "*")
		if err != nil {
			p.Log(logging.Error, "Error restoring files: %s", err)
			p.RunningEnvironment.DisplayToConsole(true, "Failed to restore files: %s", err)
		}
	}(fileName, c)

	return nil
}

func (p *Server) GetBackup(fileName string) (*FileData, error) {
	info, err := p.GetBackupFileServer().Stat(fileName)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return &FileData{ContentLength: info.Size(), Name: info.Name()}, nil
}

func (p *Server) GetBackupFile(fileName string) (*FileData, error) {
	file, err := p.GetBackupFileServer().Open(fileName)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return &FileData{Contents: file, ContentLength: info.Size(), Name: info.Name()}, nil
}

func (p *Server) valid() error {
	//we need a type at least, this is a safe check
	if p.Type.Type == "" {
		return ErrServerTypeRequired
	}

	if p.Environment.Type == "" {
		return ErrEnvironmentTypeRequired
	}

	return nil
}

func (p *Server) Log(l *log.Logger, format string, obj ...interface{}) {
	msg := fmt.Sprintf("[%s] ", p.Id()) + format
	l.Printf(msg, obj...)
}

func (p *Server) RunCondition(condition string, extraData map[string]interface{}) (bool, error) {
	data := map[string]interface{}{
		conditions.VariableEnv:      p.RunningEnvironment.Type,
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

	return conditions.ResolveIf(condition, data, CreateFunctions(p, p.GetEnvironment()))
}

func (p *Server) GetFileServer() files.FileServer {
	return p.fileServer
}

func (p *Server) GetBackupFileServer() files.FileServer {
	return p.backupFileServer
}

func (p *Server) IsBackingUp() bool {
	return p.backingUp
}

func (p *Server) IsRestoring() bool {
	return p.restoring
}

func (p *Server) IsIdle() error {
	if p.IsRestoring() || p.IsBackingUp() {
		return pufferpanel.ErrBackupInProgress
	}

	r, _ := p.GetEnvironment().IsRunning()
	if r {
		return pufferpanel.ErrServerRunning
	}

	if p.GetEnvironment().IsInstalling() {
		return pufferpanel.ErrServerRunning
	}

	return nil
}

func (p *Server) SetFileServer(fs files.FileServer) {
	p.fileServer = fs
}
