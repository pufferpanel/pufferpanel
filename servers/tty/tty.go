//go:build !windows

package tty

import (
	"errors"
	"fmt"
	"github.com/creack/pty"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/shirou/gopsutil/process"
	"github.com/spf13/cast"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

type tty struct {
	*pufferpanel.BaseEnvironment
	mainProcess *exec.Cmd

	statLocker   sync.Mutex
	lastStats    *pufferpanel.ServerStats
	lastStatTime time.Time
}

func (t *tty) ttyExecuteAsync(steps pufferpanel.ExecutionData) (err error) {
	t.Wait.Add(1)

	pr := exec.Command(steps.Command, steps.Arguments...)
	pr.Dir = t.GetRootDirectory()
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

	_ = t.StatusTracker.WriteMessage(pufferpanel.Transmission{
		Message: pufferpanel.ServerRunning{
			Running:    true,
			Installing: t.IsInstalling(),
		},
		Type: pufferpanel.MessageTypeStatus,
	})

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
		stats := &pufferpanel.ServerStats{
			Cpu:    0,
			Memory: 0,
		}

		if t.Server.Stats.Type == "jcmd" {
			stats.Jvm = &pufferpanel.JvmStats{}
		}

		return stats, nil
	}

	t.statLocker.Lock()
	defer t.statLocker.Unlock()

	//only fetch stats once every 5 seconds, to avoid excessive spam
	if t.lastStatTime.Add(5 * time.Second).After(time.Now()) {
		return t.lastStats, nil
	}

	pr, err := process.NewProcess(int32(t.mainProcess.Process.Pid))
	if err != nil {
		return nil, err
	}

	memMap, _ := pr.MemoryInfo()
	cpu, _ := pr.Percent(time.Second * 1)

	stats := &pufferpanel.ServerStats{
		Cpu:    cpu,
		Memory: cast.ToFloat64(memMap.RSS),
	}

	if t.Server.Stats.Type == "jcmd" {
		var socket *net.UnixConn
		if socket, err = t.initiateJCMD(); err == nil && socket != nil {
			for _, s := range []string{"1", "\x00", "jcmd", "\x00", "GC.heap_info", "\x00", "\x00", "\x00"} {
				_, err = socket.Write([]byte(s))
				if err != nil {
					logging.Error.Printf("unable to send command to Java process: %v", err)
					break
				}
			}
			//only continue parsing if no errors sending command
			if err == nil {
				var jcmdData []byte
				jcmdData, err = io.ReadAll(socket)
				if err != nil {
					logging.Error.Printf("Could not get result of JCMD: %s", err.Error())
				}

				stats.Jvm = pufferpanel.ParseJCMDResponse(jcmdData)
			}
		}
		if stats.Jvm == nil {
			stats.Jvm = &pufferpanel.JvmStats{}
		}
	}

	t.lastStats = stats

	return stats, nil
}

func (t *tty) SendCode(code int) error {
	running, err := t.IsRunning()

	if err != nil || !running {
		return err
	}

	return t.mainProcess.Process.Signal(syscall.Signal(code))
}

func (t *tty) GetUid() int {
	return -1
}

func (t *tty) GetGid() int {
	return -1
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

	t.statLocker.Lock()
	t.statLocker.Unlock()

	t.mainProcess = nil

	t.Wait.Done()

	_ = t.StatusTracker.WriteMessage(pufferpanel.Transmission{
		Message: pufferpanel.ServerRunning{
			Running:    false,
			Installing: t.IsInstalling(),
		},
		Type: pufferpanel.MessageTypeStatus,
	})

	if callback != nil {
		callback(exitCode)
	}
}

func activateAttachAPI(pid int) error {
	// It's not, lets do a quick ceremony of touching a file and
	// sending SIGQUIT to activate this feature
	attachpath := attachPath(pid)
	if err := os.WriteFile(attachpath, nil, 0660); err != nil {
		return fmt.Errorf("could not touch file to activate attach api: %w", err)
	}

	defer func() {
		_ = os.Remove(attachpath)
	}()

	proc, err := os.FindProcess(pid)
	if err != nil { // can't happen on unix
		return fmt.Errorf("could not find process: %w", err)
	}

	if err = proc.Signal(syscall.SIGQUIT); err != nil {
		return fmt.Errorf("could not send signal 3 to activate attach API: %w", err)
	}

	// Check if the UNIX socket is active
	sock := socketPath(pid)
	for i := 1; i < 10; i++ {
		if _, err = os.Stat(sock); err != nil && !os.IsNotExist(err) {
			return err
		}

		// exponential backoff
		time.Sleep(time.Duration(1<<uint(i)) * time.Millisecond)
	}

	//if we got here, then the file wasn't available or otherwise not good anymore
	return err
}

func attachPath(pid int) string {
	return fmt.Sprintf("/proc/%v/cwd/.attach_pid%v", pid, pid)
}

func socketPath(pid int) string {
	return fmt.Sprintf("/proc/%v/root/tmp/.java_pid%v", pid, pid)
}

func (t *tty) initiateJCMD() (*net.UnixConn, error) {
	pid := t.mainProcess.Process.Pid
	sock := socketPath(pid)

	// Check if the UNIX socket is active
	if _, err := os.Stat(sock); err != nil && os.IsNotExist(err) {
		if err = activateAttachAPI(pid); err != nil {
			return nil, err
		}
	}

	addr, err := net.ResolveUnixAddr("unix", sock)
	if err != nil {
		return nil, err // can't happen (on linux)
	}

	return net.DialUnix("unix", nil, addr)
}
