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
	//Executes a command within the environment.
	Execute(steps ExecutionData) error

	//Executes a command within the environment and immediately return
	ExecuteAsync(steps ExecutionData) error

	//Sends a string to the StdIn of the main program process
	ExecuteInMainProcess(cmd string) error

	//Kills the main process, but leaves the environment running.
	Kill() error

	//Creates the environment setting needed to run programs.
	Create() error

	//Deletes the environment.
	Delete() error

	Update() error

	IsRunning() (isRunning bool, err error)

	WaitForMainProcess() error

	WaitForMainProcessFor(timeout time.Duration) error

	GetRootDirectory() string

	GetConsole() (console []string, epoch int64)

	GetConsoleFrom(time int64) (console []string, epoch int64)

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

	SendStats()
}

type BaseEnvironment struct {
	Environment
	Type              string               `json:"type"`
	RootDirectory     string               `json:"root,omitempty"`
	ConsoleBuffer     Cache                `json:"-"`
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

func (e *BaseEnvironment) GetConsole() (console []string, epoch int64) {
	console, epoch = e.ConsoleBuffer.Read()
	return
}

func (e *BaseEnvironment) GetConsoleFrom(time int64) (console []string, epoch int64) {
	console, epoch = e.ConsoleBuffer.ReadFrom(time)
	return
}

func (e *BaseEnvironment) AddConsoleListener(ws *Socket) {
	e.ConsoleTracker.Register(ws)
}

func (e *BaseEnvironment) AddStatListener(ws *Socket) {
	e.StatsTracker.Register(ws)
}

func (e *BaseEnvironment) AddStatusListener(ws *Socket) {
	e.StatusTracker.Register(ws)
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

func (e *BaseEnvironment) SendStats() {
	stats, err := e.GetStats()
	if err != nil {
		return
	}
	_ = e.StatsTracker.WriteMessage(&messages.Stat{
		Memory: stats.Memory,
		Cpu:    stats.Cpu,
	})
}

func (e *BaseEnvironment) Log(l *log.Logger, format string, obj ...interface{}) {
	msg := fmt.Sprintf("[%s] ", e.ServerId) + format
	l.Printf(msg, obj...)
}

func newLogger(prefix string) *log.Logger {
	return log.New(logging.Info.Writer(), "["+prefix+"] ", 0)
}
