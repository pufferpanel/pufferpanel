//go:build !windows

/*
 Copyright 2020 Padduck, LLC
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

package testing

import (
	"fmt"
	"github.com/creack/pty"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
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

func CreateEnvironment(prefix string) *Environment {
	rootDir := "/var/lib/pufferpanel/testing"
	if prefix != "" {
		rootDir = filepath.Join(rootDir, prefix)
	}

	t := &Environment{
		BaseEnvironment: &pufferpanel.BaseEnvironment{Type: "tty", RootDirectory: rootDir},
	}
	t.BaseEnvironment.ExecutionFunction = t.ttyExecuteAsync
	t.BaseEnvironment.WaitFunction = t.WaitForMainProcess
	t.Wait = &sync.WaitGroup{}
	return t
}

func (t *Environment) ttyExecuteAsync(steps pufferpanel.ExecutionData) (err error) {
	running, err := t.IsRunning()
	if err != nil {
		return
	}
	if running {
		err = pufferpanel.ErrProcessRunning
		return
	}
	t.Wait.Wait()

	pr := exec.Command(steps.Command, steps.Arguments...)
	pr.Dir = path.Join(t.GetRootDirectory(), steps.WorkingDirectory)
	/*for _, v := range os.Environ() {
		if !strings.HasPrefix(v, "PUFFER_") {
			pr.Env = append(pr.Env, v)
		}
	}*/
	pr.Env = append(pr.Env, "HOME="+t.GetRootDirectory(), "TERM=xterm-256color")
	for k, v := range steps.Environment {
		pr.Env = append(pr.Env, fmt.Sprintf("%s=%s", k, v))
	}

	wrapper := t.CreateWrapper()
	t.Wait.Add(1)
	pr.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}
	t.mainProcess = pr
	t.DisplayToConsole(true, "Starting process: %s %s", t.mainProcess.Path, strings.Join(t.mainProcess.Args[1:], " "))

	tty, err := pty.Start(pr)
	if err != nil {
		return
	}

	t.stdInWriter = tty

	go func(proxy io.Writer) {
		_, _ = io.Copy(proxy, tty)
	}(wrapper)

	go t.handleClose(steps.Callback)
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

func (t *Environment) handleClose(callback func(exitCode int)) {
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
		callback(code)
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
