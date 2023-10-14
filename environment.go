package pufferpanel

import (
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/messages"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type Environment interface {
	// Execute Executes a command within the environment.
	Execute(steps ExecutionData) error

	// ExecuteAsync Executes a command within the environment and immediately return
	ExecuteAsync(steps ExecutionData) error

	// ExecuteInMainProcess Sends a string to the StdIn of the main program process
	ExecuteInMainProcess(cmd string) error

	// Kill Kills the main process, but leaves the environment running.
	Kill() error

	// Create Creates the environment setting needed to run programs.
	Create() error

	// Delete Deletes the environment.
	Delete() error

	Update() error

	IsRunning() (isRunning bool, err error)

	IsInstalling() bool

	SetInstalling(bool)

	WaitForMainProcess() error

	WaitForMainProcessFor(timeout time.Duration) error

	GetRootDirectory() string

	GetConsole() (console []byte, epoch int64)

	GetConsoleFrom(time int64) (console []byte, epoch int64)

	AddConsoleListener(ws *Socket)

	AddStatusListener(ws *Socket)

	AddStatsListener(ws *Socket)

	GetStats() (*ServerStats, error)

	DisplayToConsole(prefix bool, msg string, data ...interface{})

	SendCode(code int) error

	GetBase() *BaseEnvironment

	GetLastExitCode() int

	GetStdOutConfiguration() ConsoleConfiguration

	GetWrapper() io.Writer

	GetStatsTracker() *Tracker
}

type BaseEnvironment struct {
	Environment
	Type              string               `json:"type"`
	RootDirectory     string               `json:"root,omitempty"`
	ConsoleBuffer     *MemoryCache         `json:"-"`
	Wait              *sync.WaitGroup      `json:"-"`
	ExecutionFunction ExecutionFunction    `json:"-"`
	WaitFunction      func() (err error)   `json:"-"`
	ServerId          string               `json:"-"`
	LastExitCode      int                  `json:"-"`
	StdOutConfig      ConsoleConfiguration `json:"stdout,omitempty"`
	StdInConfig       ConsoleConfiguration `json:"stdin,omitempty"`
	Wrapper           io.Writer            `json:"-"` //our proxy back to the main
	ConsoleTracker    *Tracker             `json:"-"`
	StatusTracker     *Tracker             `json:"-"`
	StatsTracker      *Tracker             `json:"-"`
	Installing        bool                 `json:"-"`
}

type ConsoleConfiguration struct {
	Type string `json:"type"`
	File string `json:"file,omitempty"`
}

type ExecutionData struct {
	Command          string
	Arguments        []string
	Environment      map[string]string
	WorkingDirectory string
	Variables        map[string]interface{}
	Callback         func(exitCode int)
}

type ExecutionFunction func(steps ExecutionData) (err error)

func (e *BaseEnvironment) Execute(steps ExecutionData) error {
	err := e.ExecuteAsync(steps)
	if err != nil {
		return err
	}
	return e.WaitForMainProcess()
}

func (e *BaseEnvironment) WaitForMainProcess() (err error) {
	return e.WaitFunction()
}

func (e *BaseEnvironment) ExecuteAsync(steps ExecutionData) (err error) {
	return e.ExecutionFunction(steps)
}

func (e *BaseEnvironment) GetRootDirectory() string {
	return e.RootDirectory
}

func (e *BaseEnvironment) GetConsole() (console []byte, epoch int64) {
	console, epoch = e.ConsoleBuffer.Read()
	return
}

func (e *BaseEnvironment) GetConsoleFrom(time int64) (console []byte, epoch int64) {
	console, epoch = e.ConsoleBuffer.ReadFrom(time)
	return
}

func (e *BaseEnvironment) AddConsoleListener(ws *Socket) {
	e.ConsoleTracker.Register(ws)
}

func (e *BaseEnvironment) AddStatsListener(ws *Socket) {
	e.StatsTracker.Register(ws)
}

func (e *BaseEnvironment) AddStatusListener(ws *Socket) {
	e.StatusTracker.Register(ws)
}

func (e *BaseEnvironment) GetStatsTracker() *Tracker {
	return e.StatsTracker
}

func (e *BaseEnvironment) DisplayToConsole(daemon bool, msg string, data ...interface{}) {
	format := msg
	if daemon {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		format = "[DAEMON] " + format
	}
	if len(data) == 0 {
		_, _ = fmt.Fprint(e.ConsoleBuffer, format)
		_, _ = fmt.Fprint(e.ConsoleTracker, format)
	} else {
		_, _ = fmt.Fprintf(e.ConsoleBuffer, format, data...)
		_, _ = fmt.Fprintf(e.ConsoleTracker, format, data...)
	}
}

func (e *BaseEnvironment) Update() error {
	return nil
}

func (e *BaseEnvironment) Delete() (err error) {
	err = os.RemoveAll(e.RootDirectory)
	return
}

func (e *BaseEnvironment) CreateWrapper() io.Writer {
	if config.ConsoleForward.Value() {
		return io.MultiWriter(newLogger(e.ServerId).Writer(), e.ConsoleBuffer, e.ConsoleTracker)
	}
	return io.MultiWriter(e.ConsoleBuffer, e.ConsoleTracker)
}

func (e *BaseEnvironment) GetBase() *BaseEnvironment {
	return e
}

func (e *BaseEnvironment) GetLastExitCode() int {
	return e.LastExitCode
}

func (e *BaseEnvironment) GetStdOutConfiguration() ConsoleConfiguration {
	return e.StdOutConfig
}

func (e *BaseEnvironment) GetWrapper() io.Writer {
	return e.Wrapper
}

func (e *BaseEnvironment) Log(l *log.Logger, format string, obj ...interface{}) {
	msg := fmt.Sprintf("[%s] ", e.ServerId) + format
	l.Printf(msg, obj...)
}

func (e *BaseEnvironment) IsInstalling() bool {
	return e.Installing
}

func (e *BaseEnvironment) SetInstalling(flag bool) {
	e.Installing = flag
	_ = e.StatusTracker.WriteMessage(&messages.Status{Installing: flag})
}

func newLogger(prefix string) *log.Logger {
	return log.New(logging.Info.Writer(), "["+prefix+"] ", 0)
}
