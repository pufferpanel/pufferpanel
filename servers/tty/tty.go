//go:build !windows

package tty

import (
	"errors"
	"fmt"
	"github.com/creack/pty"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/messages"
	"github.com/shirou/gopsutil/process"
	"github.com/spf13/cast"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"
)

type tty struct {
	*pufferpanel.BaseEnvironment
	mainProcess *exec.Cmd
}

func (t *tty) ttyExecuteAsync(steps pufferpanel.ExecutionData) (err error) {
	t.Wait.Add(1)

	pr := exec.Command(steps.Command, steps.Arguments...)
	pr.Dir = path.Join(t.GetRootDirectory(), steps.WorkingDirectory)
	for _, v := range os.Environ() {
		if !strings.HasPrefix(v, "PUFFER_") {
			pr.Env = append(pr.Env, v)
		}
	}
	pr.Env = append(pr.Env, "HOME="+t.GetRootDirectory(), "TERM=xterm-256color")
	for k, v := range steps.Environment {
		pr.Env = append(pr.Env, fmt.Sprintf("%s=%s", k, v))
	}

	pr.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}
	t.mainProcess = pr
	t.DisplayToConsole(true, "Starting process: %s %s", t.mainProcess.Path, strings.Join(t.mainProcess.Args[1:], " "))
	t.Log(logging.Info, "Starting process: %s %s", t.mainProcess.Path, strings.Join(t.mainProcess.Args[1:], " "))

	msg := messages.Status{Running: true, Installing: t.IsInstalling()}
	_ = t.StatusTracker.WriteMessage(msg)

	processTty, err := pty.Start(pr)
	if err != nil {
		t.Wait.Done()
		return
	}

	t.BaseEnvironment.CreateConsoleStdinProxy(steps.StdInConfig, processTty)
	t.BaseEnvironment.Console.Start()

	go func(proxy io.Writer) {
		_, _ = io.Copy(proxy, processTty)
	}(t.Wrapper)

	go t.handleClose(steps.Callback)
	return
}

func (t *tty) kill() (err error) {
	running, err := t.IsRunning()
	if err != nil {
		return
	}
	if !running {
		return
	}
	return t.mainProcess.Process.Kill()
}

func (t *tty) GetStats() (*pufferpanel.ServerStats, error) {
	running, err := t.IsRunning()
	if err != nil {
		return nil, err
	}
	if !running {
		return &pufferpanel.ServerStats{
			Cpu:    0,
			Memory: 0,
		}, nil
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

func (t *tty) SendCode(code int) error {
	running, err := t.IsRunning()

	if err != nil || !running {
		return err
	}

	return t.mainProcess.Process.Signal(syscall.Signal(code))
}

func (t *tty) isRunning() (isRunning bool, err error) {
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

func (t *tty) handleClose(callback func(exitCode int)) {
	err := t.mainProcess.Wait()

	_ = t.Console.Close()

	var exitCode int
	if t.mainProcess.ProcessState == nil || err != nil {
		var psErr *exec.ExitError
		if errors.As(err, &psErr) {
			exitCode = psErr.ExitCode()
		} else {
			exitCode = 1
		}
	} else {
		exitCode = t.mainProcess.ProcessState.ExitCode()
	}
	t.LastExitCode = exitCode

	if err != nil {
		t.Log(logging.Error, "Error waiting on process: %s\n", err)
	}

	if t.mainProcess != nil && t.mainProcess.ProcessState != nil {
		t.Log(logging.Debug, "%s\n", t.mainProcess.ProcessState.String())
	}

	if t.mainProcess != nil && t.mainProcess.Process != nil {
		_ = t.mainProcess.Process.Release()
	}
	t.mainProcess = nil
	t.Wait.Done()

	msg := messages.Status{Running: false, Installing: t.IsInstalling()}
	_ = t.StatusTracker.WriteMessage(msg)

	if callback != nil {
		callback(exitCode)
	}
}
