//go:build windows
// +build windows

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

package test

import (
	"fmt"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/messages"
	"github.com/shirou/gopsutil/process"
	"github.com/spf13/cast"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Environment struct {
	*pufferpanel.BaseEnvironment
	mainProcess *exec.Cmd
	stdInWriter io.Writer
}

func CreateEnvironment() *Environment {
	t := &Environment{
		BaseEnvironment: &pufferpanel.BaseEnvironment{Type: "standard"},
	}
	t.BaseEnvironment.ExecutionFunction = t.standardExecuteAsync
	t.BaseEnvironment.WaitFunction = t.WaitForMainProcess
	t.Wait = &sync.WaitGroup{}
	return t
}

func (s *Environment) standardExecuteAsync(steps pufferpanel.ExecutionData) (err error) {
	running, err := s.IsRunning()
	if err != nil {
		return
	}
	if running {
		err = pufferpanel.ErrProcessRunning
		return
	}
	s.Wait.Wait()
	s.Wait.Add(1)
	s.mainProcess = exec.Command(steps.Command, steps.Arguments...)
	s.mainProcess.Dir = path.Join(s.GetRootDirectory(), steps.WorkingDirectory)

	for _, v := range os.Environ() {
		if !strings.HasPrefix(v, "PUFFER_") {
			s.mainProcess.Env = append(s.mainProcess.Env, v)
		}
	}
	s.mainProcess.Env = append(s.mainProcess.Env, "HOME="+s.GetRootDirectory(), "TERM=xterm-256color")
	for k, v := range steps.Environment {
		s.mainProcess.Env = append(s.mainProcess.Env, fmt.Sprintf("%s=%s", k, v))
	}
	wrapper := s.CreateWrapper()
	s.mainProcess.Stdout = wrapper
	s.mainProcess.Stderr = wrapper
	pipe, err := s.mainProcess.StdinPipe()
	if err != nil {
		return
	}
	s.stdInWriter = pipe
	logging.Info.Printf("Starting process: %s %s", s.mainProcess.Path, strings.Join(s.mainProcess.Args[1:], " "))
	s.DisplayToConsole(true, "Starting process: %s %s", s.mainProcess.Path, strings.Join(s.mainProcess.Args[1:], " "))

	msg := messages.Status{Running: true}
	_ = s.WSManager.WriteMessage(msg)

	err = s.mainProcess.Start()
	if err != nil && err.Error() != "exit status 1" {
		msg := messages.Status{Running: false}
		_ = s.WSManager.WriteMessage(msg)
		logging.Info.Printf("Process failed to start: %s", err)
		return
	} else {
		logging.Info.Printf("Process started (%d)", s.mainProcess.Process.Pid)
	}

	go s.handleClose(steps.Callback)
	return
}

func (t *Environment) ExecuteInMainProcess(cmd string) (err error) {
	running, err := t.IsRunning()
	if err != nil {
		return err
	}
	if !running {
		err = pufferpanel.ErrServerOffline
		return
	}
	stdIn := t.stdInWriter
	_, err = io.WriteString(stdIn, cmd+"\n")
	return
}

func (t *Environment) Kill() (err error) {
	running, err := t.IsRunning()
	if err != nil {
		return
	}
	if !running {
		return
	}
	return t.mainProcess.Process.Kill()
}

func (t *Environment) IsRunning() (isRunning bool, err error) {
	isRunning = t.mainProcess != nil && t.mainProcess.Process != nil
	if isRunning {
		pr, pErr := os.FindProcess(t.mainProcess.Process.Pid)
		if pr == nil || pErr != nil {
			isRunning = false
		} else if pr.Signal(syscall.Signal(0)) != nil {
			isRunning = false
		}
	}
	return
}

func (t *Environment) GetStats() (*pufferpanel.ServerStats, error) {
	running, err := t.IsRunning()
	if err != nil {
		return nil, err
	}
	if !running {
		return nil, pufferpanel.ErrServerOffline
	}
	pr, err := process.NewProcess(int32(t.mainProcess.Process.Pid))
	if err != nil {
		return nil, err
	}

	memMap, _ := pr.MemoryInfo()
	cpu, _ := pr.Percent(time.Second * 1)

	return &pufferpanel.ServerStats{
		Cpu:    cpu,
		Memory: cast.ToFloat64(memMap.RSS),
	}, nil
}

func (t *Environment) Create() error {
	return os.Mkdir(t.GetRootDirectory(), 0755)
}

func (t *Environment) WaitForMainProcess() error {
	return t.WaitForMainProcessFor(0)
}

func (t *Environment) WaitForMainProcessFor(timeout time.Duration) (err error) {
	running, err := t.IsRunning()
	if err != nil {
		return
	}
	if running {
		if timeout > 0 {
			var timer = time.AfterFunc(timeout, func() {
				err = t.Kill()
			})
			t.Wait.Wait()
			timer.Stop()
		} else {
			t.Wait.Wait()
		}
	}
	return
}

func (t *Environment) SendCode(code int) error {
	running, err := t.IsRunning()

	if err != nil || !running {
		return err
	}

	return t.mainProcess.Process.Signal(syscall.Signal(code))
}

func (t *Environment) handleClose(callback func(exitCode bool)) {
	err := t.mainProcess.Wait()

	var code int
	if t.mainProcess == nil || t.mainProcess.ProcessState == nil || err != nil {
		code = 1
	} else {
		code = t.mainProcess.ProcessState.ExitCode()
	}

	if t.mainProcess != nil && t.mainProcess.Process != nil {
		_ = t.mainProcess.Process.Release()
	}
	t.mainProcess = nil
	t.Wait.Done()

	if callback != nil {
		callback(code == 0)
	}
}

func (*Environment) CreateWrapper() io.Writer {
	return os.Stdout
}

func (*Environment) DisplayToConsole(prefix bool, msg string, data ...interface{}) {
	if prefix {
		logging.Info.Printf("[DAEMON] "+msg, data...)
	} else {
		logging.Info.Printf(msg, data...)
	}
}

func (*Environment) GetRootDirectory() string {
	return "C:\\Temp\\PufferPanel\\testing"
}
