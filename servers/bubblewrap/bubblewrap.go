package bubblewrap

import (
	"errors"
	"fmt"
	"github.com/creack/pty"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/utils"
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

type bubblewrap struct {
	mainProcess  *exec.Cmd
	statLocker   sync.Mutex
	lastStats    *pufferpanel.ServerStats
	lastStatTime time.Time
}

func (t *bubblewrap) ExecuteAsyncImpl(environment *pufferpanel.Environment, steps pufferpanel.ExecutionData) (err error) {
	environment.Wait.Add(1)

	pr := exec.Command(steps.Command, steps.Arguments...)
	pr.Dir = environment.GetRootDirectory()

	var envVars = make(map[string]string)

	for _, v := range os.Environ() {
		key, value, valid := strings.Cut(v, "=")
		if !valid {
			continue
		}
		if strings.HasPrefix(key, "PUFFER_") {
			continue
		}
		envVars[key] = value
	}
	envVars["HOME"] = environment.GetRootDirectory()
	envVars["TERM"] = "xterm-256color"
	for k, v := range steps.Environment {
		envVars[k] = v
	}

	for k, v := range envVars {
		pr.Env = append(pr.Env, fmt.Sprintf("%s=%s", k, v))
	}

	pr.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}
	t.mainProcess = pr
	environment.DisplayToConsole(true, "Starting process: %s %s", t.mainProcess.Path, strings.Join(t.mainProcess.Args[1:], " "))
	environment.Log(logging.Info, "Starting process: %s %s", t.mainProcess.Path, strings.Join(t.mainProcess.Args[1:], " "))

	_ = environment.StatusTracker.WriteMessage(pufferpanel.Transmission{
		Message: pufferpanel.ServerRunning{
			Running:    true,
			Installing: environment.IsInstalling(),
		},
		Type: pufferpanel.MessageTypeStatus,
	})

	processTty, err := pty.Start(pr)
	if err != nil {
		environment.Wait.Done()
		return
	}

	environment.CreateConsoleStdinProxy(steps.StdInConfig, processTty)
	environment.Console.Start()

	go func(proxy io.Writer) {
		_, _ = io.Copy(proxy, processTty)
	}(environment.Wrapper)

	go t.handleClose(environment, steps.Callback)
	return
}

func (t *bubblewrap) KillImpl(environment *pufferpanel.Environment) (err error) {
	running, err := environment.IsRunning()
	if err != nil {
		return
	}
	if !running {
		return
	}
	return t.mainProcess.Process.Kill()
}

func (t *bubblewrap) GetStatsImpl(environment *pufferpanel.Environment) (*pufferpanel.ServerStats, error) {
	running, err := environment.IsRunning()
	if err != nil {
		return nil, err
	}
	if !running {
		stats := &pufferpanel.ServerStats{
			Cpu:    0,
			Memory: 0,
		}

		if environment.Server.Stats.Type == "jcmd" {
			stats.Jvm = &utils.JvmStats{}
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

	if environment.Server.Stats.Type == "jcmd" {
		if socket, err := t.initializeJCmd(); err == nil && socket != nil {
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

				stats.Jvm = utils.ParseJCMDResponse(jcmdData)
			}
		}
		if stats.Jvm == nil {
			stats.Jvm = &utils.JvmStats{}
		}
	}

	t.lastStats = stats

	return stats, nil
}

func (t *bubblewrap) SendCodeImpl(environment *pufferpanel.Environment, code int) error {
	running, err := environment.IsRunning()

	if err != nil || !running {
		return err
	}

	return t.mainProcess.Process.Signal(syscall.Signal(code))
}

func (t *bubblewrap) GetUidImpl(*pufferpanel.Environment) int {
	return -1
}

func (t *bubblewrap) GetGidImpl(*pufferpanel.Environment) int {
	return -1
}

func (t *bubblewrap) IsRunningImpl(*pufferpanel.Environment) (isRunning bool, err error) {
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

func (t *bubblewrap) handleClose(environment *pufferpanel.Environment, callback func(exitCode int)) {
	err := t.mainProcess.Wait()

	_ = environment.Console.Close()

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
	environment.LastExitCode = exitCode

	if err != nil {
		environment.Log(logging.Error, "Error waiting on process: %s\n", err)
	}

	if t.mainProcess != nil && t.mainProcess.ProcessState != nil {
		environment.Log(logging.Debug, "%s\n", t.mainProcess.ProcessState.String())
	}

	if t.mainProcess != nil && t.mainProcess.Process != nil {
		_ = t.mainProcess.Process.Release()
	}

	t.statLocker.Lock()
	t.statLocker.Unlock()

	t.mainProcess = nil

	environment.Wait.Done()

	_ = environment.StatusTracker.WriteMessage(pufferpanel.Transmission{
		Message: pufferpanel.ServerRunning{
			Running:    false,
			Installing: environment.IsInstalling(),
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

func (t *bubblewrap) initializeJCmd() (net.Conn, error) {
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
