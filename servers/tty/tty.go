package tty

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/creack/pty"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/utils"
	"github.com/shirou/gopsutil/process"
	"github.com/spf13/cast"
)

type tty struct {
	mainProcess  *exec.Cmd
	statLocker   sync.Mutex
	lastStats    *pufferpanel.ServerStats
	lastStatTime time.Time
	//disableStdin        bool
	disableSpecialStats bool

	DisableUnshare bool     `json:"disableUnshare"`
	Mounts         []string `json:"mounts"`
}

func (t *tty) ExecuteAsyncImpl(environment *pufferpanel.Environment, steps pufferpanel.ExecutionData) (err error) {
	pr, err := t.createCmd(environment.GetRootDirectory(), steps.Command)
	if err != nil {
		return err
	}

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
	envVars["PUFFERPANEL_SERVER_COMMAND"] = steps.Command
	for k, v := range steps.Environment {
		envVars[k] = v
	}

	for k, v := range envVars {
		pr.Env = append(pr.Env, fmt.Sprintf("%s=%s", k, v))
	}

	t.mainProcess = pr
	environment.DisplayToConsole(true, "Starting process: %s", steps.Command)
	environment.Log(logging.Info, "Starting process in directory [%s]: %s", t.mainProcess.Dir, strings.Replace(strings.Join(t.mainProcess.Args, " "), "${PUFFERPANEL_SERVER_COMMAND}", steps.Command, 1))

	_ = environment.StatusTracker.WriteMessage(pufferpanel.Transmission{
		Message: pufferpanel.ServerRunning{
			Running:    true,
			Installing: environment.IsInstalling(),
		},
		Type: pufferpanel.MessageTypeStatus,
	})

	t.disableSpecialStats = steps.DisableStats
	//t.disableStdin = steps.DisableStdin

	processTty, err := pty.Start(pr)
	if err != nil {
		return
	}

	//if !t.disableStdin {
	//	environment.CreateConsoleStdinProxy(steps.StdInConfig, processTty)
	//}
	environment.CreateConsoleStdinProxy(steps.StdInConfig, processTty)

	environment.Console.Start()

	go func(proxy io.Writer) {
		_, _ = io.Copy(proxy, processTty)
	}(environment.Wrapper)

	go t.handleClose(environment, steps.Callback)
	return
}

func (t *tty) KillImpl(environment *pufferpanel.Environment) (err error) {
	running, err := environment.IsRunning()
	if err != nil {
		return
	}
	if !running {
		return
	}
	return t.mainProcess.Process.Kill()
}

func (t *tty) GetStatsImpl(environment *pufferpanel.Environment) (*pufferpanel.ServerStats, error) {
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

	if !t.disableSpecialStats && environment.Server.Stats.Type == "jcmd" {
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

func (t *tty) SendCodeImpl(environment *pufferpanel.Environment, code int) error {
	running, err := environment.IsRunning()

	if err != nil || !running {
		return err
	}

	return t.mainProcess.Process.Signal(syscall.Signal(code))
}

func (t *tty) GetUidImpl(*pufferpanel.Environment) int {
	return -1
}

func (t *tty) GetGidImpl(*pufferpanel.Environment) int {
	return -1
}

func (t *tty) IsRunningImpl(*pufferpanel.Environment) (isRunning bool, err error) {
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

func (t *tty) handleClose(environment *pufferpanel.Environment, callback func(exitCode int)) {
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

	//if we are using unshare AND we're in tmp, we can nuke the workspace at this point
	if !t.DisableUnshare && strings.HasPrefix(t.mainProcess.Dir, os.TempDir()) {
		err = os.RemoveAll(t.mainProcess.Dir)
		if err != nil {
			logging.Debug.Printf("Failed to delete %s: %s", t.mainProcess.Dir, err.Error())
		}
	}

	t.mainProcess = nil

	_ = environment.StatusTracker.WriteMessage(pufferpanel.Transmission{
		Message: pufferpanel.ServerRunning{
			Running:    false,
			Installing: environment.IsInstalling(),
		},
		Type: pufferpanel.MessageTypeStatus,
	})

	//t.disableStdin = false
	t.disableSpecialStats = false

	environment.Wait.Unlock()

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

var cmdList = []string{
	"mount --make-rprivate --make-rslave --bind . .",
	"mkdir -p dev bin usr lib etc tmp proc",
	"mount -t tmpfs -o size=50m tmpfs tmp",
	"mount --bind /bin bin",
	"mount --bind /lib lib",
	"mount --rbind /usr usr",
	"mount --rbind /etc etc",
	"mount --rbind /dev dev",
	"mount --rbind /proc proc",
}

const Shell = "bash"

func (t *tty) createCmd(workDir, cmd string) (pr *exec.Cmd, err error) {
	if t.DisableUnshare || config.SecurityDisableUnshare.Value() {
		pr = exec.Command(Shell, "-c", "${PUFFERPANEL_SERVER_COMMAND}")
		pr.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}
		pr.Dir = workDir
		return
	}

	workDirMount := removeRoot(workDir)
	binaryFolderMount := removeRoot(config.BinariesFolder.Value())
	cacheFolderMount := removeRoot(config.CacheFolder.Value())

	mountFolders := []string{workDirMount, binaryFolderMount, cacheFolderMount}
	for _, v := range t.Mounts {
		mountFolders = append(mountFolders, removeRoot(v))
	}

	unshareArgs := make([]string, len(cmdList))
	copy(unshareArgs, cmdList)

	if runtime.GOARCH == "amd64" {
		unshareArgs = append(unshareArgs,
			"mkdir -p lib64",
			"mount --bind /lib64 lib64",
		)
	}

	var lstat os.FileInfo
	lstat, err = os.Lstat("/etc/resolv.conf")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return
	}
	if err == nil && lstat.Mode()&os.ModeSymlink != 0 {
		var absPath string
		absPath, err = filepath.EvalSymlinks("/etc/resolv.conf")
		if err != nil {
			return
		}
		localPath := removeRoot(absPath)
		dir := removeRoot(filepath.Dir(absPath))
		unshareArgs = append(unshareArgs,
			fmt.Sprintf("mkdir -p %s", dir),
			fmt.Sprintf("touch %s", localPath),
			fmt.Sprintf("mount --rbind %s %s", absPath, localPath),
		)
	}

	unshareArgs = append(unshareArgs,
		fmt.Sprintf("mkdir -p %s", strings.Join(mountFolders, " ")),
		fmt.Sprintf("mount --bind %s %s", workDir, workDirMount),
		fmt.Sprintf("mount --bind %s %s", config.BinariesFolder.Value(), binaryFolderMount),
		fmt.Sprintf("mount --bind %s %s", config.CacheFolder.Value(), cacheFolderMount),
	)

	for _, v := range t.Mounts {
		unshareArgs = append(unshareArgs, fmt.Sprintf("mount --bind %s %s", v, removeRoot(v)))
	}

	unshareArgs = append(unshareArgs,
		//move cwd to bind mounted instace of .
		"cd .",
		"mkdir -p old-root",
		//make . the root for everything in the current namespace
		"pivot_root . old-root",
		//make the old root unaccessible by unmounting it
		//needs to be lazy because the old root is considered busy as it's still the root outside the namespace
		"umount -l /old-root",
		"rm -r /old-root",
		fmt.Sprintf("unshare -U -w %s --map-user=%d --map-group=%d %s -c '${PUFFERPANEL_SERVER_COMMAND}'", workDir, os.Getuid(), os.Getgid(), Shell))

	pr = exec.Command(Shell, "-c", strings.Join(unshareArgs, " && "))
	pr.Dir, err = os.MkdirTemp("", "unshare-pp-")
	if err != nil {
		return
	}
	pr.SysProcAttr = &syscall.SysProcAttr{
		Setctty: true,
		Setsid:  true,
		Unshareflags: syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWNS |
			syscall.CLONE_FILES |
			syscall.CLONE_NEWCGROUP |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWUTS,
		Credential: &syscall.Credential{Uid: 0, Gid: 0, NoSetGroups: true},
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}
	return
}

func removeRoot(path string) string {
	return strings.TrimPrefix(path, "/")
}
