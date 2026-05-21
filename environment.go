package pufferpanel

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/connections"
	"github.com/pufferpanel/pufferpanel/v3/files"
	"github.com/pufferpanel/pufferpanel/v3/logging"
)

type EnvironmentImpl interface {
	ExecuteAsyncImpl(environment *Environment, steps ExecutionData) error

	KillImpl(environment *Environment) error

	GetStatsImpl(environment *Environment) (*ServerStats, error)

	SendCodeImpl(environment *Environment, code int) error

	GetUidImpl(environment *Environment) int

	GetGidImpl(environment *Environment) int

	IsRunningImpl(environment *Environment) (isRunning bool, err error)
}

type Environment struct {
	Type            string          `json:"type"`
	RootDirectory   string          `json:"root,omitempty"`
	BackupDirectory string          `json:"-"`
	ConsoleBuffer   *MemoryCache    `json:"-"`
	Wait            *sync.Mutex     `json:"-"`
	ServerId        string          `json:"-"`
	LastExitCode    int             `json:"-"`
	Wrapper         io.Writer       `json:"-"` //our proxy back to the main
	ConsoleTracker  *Tracker        `json:"-"`
	StatusTracker   *Tracker        `json:"-"`
	StatsTracker    *Tracker        `json:"-"`
	Installing      bool            `json:"-"`
	BackingUp       bool            `json:"-"`
	Console         Console         `json:"-"`
	Server          Server          `json:"-"`
	Implementation  EnvironmentImpl `json:"-"`
}

type ExecutionData struct {
	Command          string
	Environment      map[string]string
	WorkingDirectory string
	Variables        map[string]interface{}
	Callback         func(exitCode int)
	StdInConfig      StdinConsoleConfiguration
	//DisableStdin     bool
	DisableQuery bool
	DisableStats bool
}

type ExecutionFunction func(steps ExecutionData) (err error)

func (e *Environment) Execute(steps ExecutionData) error {
	err := e.ExecuteAsync(steps)
	if err != nil {
		return err
	}
	return e.WaitForMainProcess()
}

func (e *Environment) ExecuteAsync(steps ExecutionData) (err error) {
	running, err := e.IsRunning()
	if err != nil {
		return
	}
	if running {
		err = ErrProcessRunning
		return
	}

	//update configs
	steps.StdInConfig = steps.StdInConfig.Replace(steps.Variables)

	e.Wait.Lock()
	err = e.Implementation.ExecuteAsyncImpl(e, steps)
	if err != nil {
		e.Wait.Unlock()
	}
	return err
}

func (e *Environment) CreateConsoleStdinProxy(config StdinConsoleConfiguration, base io.WriteCloser) {
	if config.Type == "telnet" {
		e.Console = &connections.TelnetConnection{
			IP:       config.IP,
			Port:     config.Port,
			Password: config.Password,
		}
	} else if config.Type == "rcon" {
		e.Console = &connections.RCONConnection{
			IP:       config.IP,
			Port:     config.Port,
			Password: config.Password,
		}
	} else if config.Type == "rconws" {
		e.Console = &connections.RCONWSConnection{
			IP:       config.IP,
			Port:     config.Port,
			Password: config.Password,
			//Environment: e,
		}
	} else {
		e.Console = &NoStartConsole{Base: base}
	}
}

func (e *Environment) GetRootDirectory() string {
	return e.RootDirectory
}

func (e *Environment) GetConsole() (console []byte, epoch int64) {
	console, epoch = e.ConsoleBuffer.Read()
	return
}

func (e *Environment) GetConsoleFrom(time int64) (console []byte, epoch int64) {
	console, epoch = e.ConsoleBuffer.ReadFrom(time)
	return
}

func (e *Environment) AddConsoleListener(ws *Socket) {
	e.ConsoleTracker.Register(ws)
}

func (e *Environment) AddStatsListener(ws *Socket) {
	e.StatsTracker.Register(ws)
}

func (e *Environment) AddStatusListener(ws *Socket) {
	e.StatusTracker.Register(ws)
}

func (e *Environment) GetStatsTracker() *Tracker {
	return e.StatsTracker
}

func (e *Environment) DisplayToConsole(daemon bool, msg string, data ...interface{}) {
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

func (e *Environment) Update() error {
	return nil
}

func (e *Environment) Delete() (err error) {
	err = files.RootServerFS.RemoveAll(e.RootDirectory)
	return
}

func (e *Environment) Create() error {
	err := files.RootServerFS.Mkdir(e.RootDirectory, 0755)
	if os.IsExist(err) {
		return nil
	}
	return err
}

func (e *Environment) WaitForMainProcess() error {
	return e.WaitForMainProcessFor(0)
}

func (e *Environment) WaitForMainProcessFor(timeout time.Duration) (err error) {
	running, err := e.IsRunning()
	if err != nil {
		return
	}
	if running {
		//if we can get the lock, the process is fully done
		//but we then need to unlock it, so we can let whatever needs to run "run"
		defer e.Wait.Unlock()
		if timeout > 0 {
			var timer = time.AfterFunc(timeout, func() {
				err = e.Kill()
			})
			e.Wait.Lock()
			timer.Stop()
		} else {
			e.Wait.Lock()
		}
	}
	return
}

func (e *Environment) CreateWrapper() {
	if config.ConsoleForward.Value() {
		//return io.MultiWriter(newLogger(e.ServerId).Writer(), e.ConsoleBuffer, e.ConsoleTracker)
		e.Wrapper = io.MultiWriter(logging.OriginalStdOut, e.ConsoleBuffer, e.ConsoleTracker)
	} else {
		e.Wrapper = io.MultiWriter(e.ConsoleBuffer, e.ConsoleTracker)
	}
}

func (e *Environment) GetLastExitCode() int {
	return e.LastExitCode
}

func (e *Environment) GetWrapper() io.Writer {
	return e.Wrapper
}

func (e *Environment) Log(l *log.Logger, format string, obj ...interface{}) {
	msg := fmt.Sprintf("[%s] ", e.ServerId) + format
	l.Printf(msg, obj...)
}

func (e *Environment) IsInstalling() bool {
	return e.Installing
}

func (e *Environment) SetInstalling(flag bool) {
	e.Installing = flag
	_ = e.StatusTracker.WriteMessage(Transmission{
		Message: ServerRunning{
			Installing: flag,
		},
		Type: MessageTypeStatus,
	})
}

func (e *Environment) ExecuteInMainProcess(cmd string) (err error) {
	running, err := e.IsRunning()
	if err != nil {
		return err
	}
	if !running {
		err = ErrServerOffline
		return
	}
	_, err = io.WriteString(e.Console, cmd+"\n")
	return
}

func (e *Environment) IsRunning() (isRunning bool, err error) {
	return e.Implementation.IsRunningImpl(e)
}

func (e *Environment) Kill() error {
	return e.Implementation.KillImpl(e)
}

func (e *Environment) GetStats() (*ServerStats, error) {
	return e.Implementation.GetStatsImpl(e)
}

func (e *Environment) SendCode(code int) error {
	return e.Implementation.SendCodeImpl(e, code)
}

func (e *Environment) GetUid() int {
	return e.Implementation.GetUidImpl(e)
}

func (e *Environment) GetGid() int {
	return e.Implementation.GetGidImpl(e)
}
