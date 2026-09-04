package servers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pufferpanel/pufferpanel/v3/utils"

	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/spf13/cast"
)

type Docker struct {
	ImageName     string               `json:"image"`
	Binds         map[string]string    `json:"bindings,omitempty"`
	Network       string               `json:"networkName,omitempty"`
	Ports         []string             `json:"portBindings,omitempty"`
	ContainerRoot string               `json:"containerRoot,omitempty"`
	HostConfig    container.HostConfig `json:"hostConfig,omitempty"`
	Labels        map[string]string    `json:"labels,omitempty"`
	Config        container.Config     `json:"config,omitempty"`

	connection   types.HijackedResponse
	cli          *client.Client
	statLocker   sync.Mutex
	lastStats    *pufferpanel.ServerStats
	lastStatTime time.Time
	//disableStdin        bool
	disableSpecialStats bool
}

func CreateDockerEnvironment() pufferpanel.EnvironmentImpl {
	return &Docker{
		ImageName: "pufferpanel/generic",
		Network:   "host",
		Ports:     make([]string, 0),
		Binds:     make(map[string]string),
		Labels:    make(map[string]string),
	}
}

func (d *Docker) ExecuteAsyncImpl(environment *pufferpanel.Environment, steps pufferpanel.ExecutionData) error {
	var err error
	var dockerClient *client.Client
	dockerClient, err = pufferpanel.GetDockerClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	//TODO: This logic may not work anymore, it's complicated to use an existing container with install/uninstall
	exists, err := pufferpanel.DoesContainerExist(environment.Server.Id(), ctx)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("docker container already exists")
	}

	err = d.createContainer(environment, steps, ctx)
	if err != nil {
		return err
	}

	d.disableSpecialStats = steps.DisableStats
	//d.disableStdin = steps.DisableStdin

	cfg := container.AttachOptions{
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Stream: true,
	}

	d.connection, err = dockerClient.ContainerAttach(ctx, environment.Server.Id(), cfg)
	if err != nil {
		return err
	}

	go func() {
		defer d.connection.Close()
		_, _ = io.Copy(environment.Wrapper, d.connection.Reader)
	}()

	//if !d.disableStdin {
	//	environment.CreateConsoleStdinProxy(steps.StdInConfig, d.connection.Conn)
	//}
	environment.CreateConsoleStdinProxy(steps.StdInConfig, d.connection.Conn)

	environment.Console.Start()

	startOpts := container.StartOptions{}

	_ = environment.StatusTracker.WriteMessage(pufferpanel.Transmission{
		Message: pufferpanel.ServerRunning{
			Running:    true,
			Installing: environment.IsInstalling(),
		},
		Type: pufferpanel.MessageTypeStatus,
	})

	environment.DisplayToConsole(true, "Starting container\n")
	err = dockerClient.ContainerStart(ctx, environment.Server.Id(), startOpts)
	if err != nil {
		return err
	}

	go d.handleClose(environment, dockerClient, steps.Callback)
	return err
}

func (d *Docker) KillImpl(environment *pufferpanel.Environment) error {
	running, err := environment.IsRunning()
	if err != nil {
		return err
	}

	if !running {
		return nil
	}

	dockerClient, err := pufferpanel.GetDockerClient()
	if err != nil {
		return err
	}
	err = dockerClient.ContainerKill(context.Background(), environment.Server.Id(), "SIGKILL")
	return err
}

func (d *Docker) IsRunningImpl(environment *pufferpanel.Environment) (bool, error) {
	dockerClient, err := pufferpanel.GetDockerClient()
	if err != nil {
		return false, err
	}

	ctx := context.Background()

	exists, err := pufferpanel.DoesContainerExist(environment.Server.Id(), ctx)
	if !exists {
		return false, err
	}

	stats, err := dockerClient.ContainerInspect(ctx, environment.Server.Id())
	if err != nil {
		return false, err
	}
	return stats.State.Running, nil
}

func (d *Docker) GetStatsImpl(environment *pufferpanel.Environment) (*pufferpanel.ServerStats, error) {
	running, err := environment.IsRunning()
	if err != nil {
		return nil, err
	}

	if !running {
		stats := &pufferpanel.ServerStats{
			Cpu:    0,
			Memory: 0,
		}

		if environment.Server.Get().Stats.Type == "jcmd" {
			stats.Jvm = &utils.JvmStats{}
		}

		return stats, nil
	}

	d.statLocker.Lock()
	defer d.statLocker.Unlock()

	//only fetch stats once every 5 seconds, to avoid excessive spam
	if d.lastStatTime.Add(5 * time.Second).After(time.Now()) {
		return d.lastStats, nil
	}

	dockerClient, err := pufferpanel.GetDockerClient()

	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	res, err := dockerClient.ContainerStats(ctx, environment.Server.Id(), false)
	defer func() {
		if res.Body != nil {
			utils.Close(res.Body)
		}
	}()
	if err != nil {
		return nil, err
	}

	data := &container.StatsResponse{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	//for java, we can get some extra data from the jcmd command
	//as such, we'll see if we can

	stats := &pufferpanel.ServerStats{
		Memory: utils.CalculateDockerMemoryPercent(data),
		Cpu:    utils.CalculateDockerCPUPercent(data),
	}

	if !d.disableSpecialStats && environment.Server.Get().Stats.Type == "jcmd" {
		cmd, _ := environment.Server.Get().Stats.Metadata["cmd"].(string)
		if cmd == "" {
			cmd = "jcmd"
		}

		r, e := dockerClient.ContainerExecCreate(context.Background(), environment.Server.Id(), container.ExecOptions{
			AttachStderr: true,
			AttachStdout: true,
			Cmd:          []string{cmd, "1", "GC.heap_info"},
		})

		if e == nil {
			rw, e := dockerClient.ContainerExecAttach(context.Background(), r.ID, container.ExecAttachOptions{
				Detach: false,
				Tty:    false,
			})
			if e != nil {
				logging.Error.Printf("Could not exec JCMD: %s", e.Error())
			} else {
				defer func(z types.HijackedResponse) {
					z.Close()
				}(rw)

				jcmdData, err := io.ReadAll(rw.Reader)
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

	d.lastStats = stats
	d.lastStatTime = time.Now()

	return stats, nil
}

func (d *Docker) createContainer(environment *pufferpanel.Environment, data pufferpanel.ExecutionData, ctx context.Context) error {
	environment.Log(logging.Debug, "Creating container")
	containerRoot := d.ContainerRoot
	if containerRoot == "" {
		containerRoot = "/pufferpanel"
	}

	if runtime.GOOS != "windows" {
		if !filepath.IsAbs(containerRoot) {
			return pufferpanel.ErrPathNotAbs(containerRoot)
		}
	}

	imageName := utils.ReplaceTokens(d.ImageName, data.Variables, utils.PlainReplace)

	err := pufferpanel.PullDockerImage(environment, ctx, imageName, false)

	if err != nil {
		return err
	}

	cmd, args := utils.SplitArguments(data.Command)

	cmdSlice := strslice.StrSlice{}
	if data.Command != "" {
		cmdSlice = append(cmdSlice, cmd)
		cmdSlice = append(cmdSlice, args...)
	}

	environment.Log(logging.Debug, "Container command: %s\n", cmdSlice)

	labels := map[string]string{
		"pufferpanel.server": environment.Server.Id(),
	}

	for k, v := range d.Labels {
		labels[utils.ReplaceTokens(k, data.Variables, utils.PlainReplace)] = utils.ReplaceTokens(v, data.Variables, utils.PlainReplace)
	}

	c := d.Config
	containerConfig := &c

	//these we need to override
	containerConfig.AttachStderr = true
	containerConfig.AttachStdin = true
	containerConfig.AttachStdout = true
	containerConfig.Tty = true
	containerConfig.OpenStdin = true
	containerConfig.NetworkDisabled = false
	containerConfig.Labels = labels

	//default if it wasn't overridden
	if containerConfig.Image == "" {
		containerConfig.Image = imageName
	}

	if containerConfig.WorkingDir == "" {
		containerConfig.WorkingDir = containerRoot
	}

	//append anything the container config added
	var envVars = make(map[string]string)

	for _, v := range containerConfig.Env {
		key, value, valid := strings.Cut(v, "=")
		if !valid {
			continue
		}
		if strings.HasPrefix(key, "PUFFER_") {
			continue
		}
		envVars[key] = value
	}
	envVars["HOME"] = containerRoot
	envVars["TERM"] = "xterm-256color"

	for k, v := range data.Environment {
		envVars[k] = v
	}

	containerConfig.Env = make([]string, 0)
	for k, v := range envVars {
		containerConfig.Env = append(containerConfig.Env, fmt.Sprintf("%s=%s", k, utils.ReplaceTokens(v, data.Variables, utils.PlainReplace)))
	}

	if len(containerConfig.Entrypoint) == 0 && len(cmdSlice) > 0 {
		containerConfig.Entrypoint = cmdSlice
	}

	if containerConfig.User == "" && runtime.GOOS != "windows" {
		containerConfig.User = fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid())
	}

	var dir string
	if containerMountSource != "" {
		dir = filepath.Join(containerMountSource, "servers", environment.Server.Id())
	} else {
		dir = filepath.Join(environment.Server.GetFileServer().Prefix(), environment.Server.Id())
	}

	//convert root dir to a full path, so we can bind it
	if !filepath.IsAbs(dir) {
		dir, err = filepath.Abs(dir)
		if err != nil {
			return err
		}
	}

	bindDirs := []string{utils.ConvertToDockerBind(dir) + ":" + containerRoot, "/etc/timezone:/etc/timezone:ro"}

	binaryFolder := config.BinariesFolder.Value()
	if containerMountSource != "" {
		binaryFolder = filepath.Join(containerMountSource, "binaries")
	} else {
		if !filepath.IsAbs(binaryFolder) {
			var ef error
			binaryFolder, ef = filepath.Abs(binaryFolder)
			if ef != nil {
				logging.Error.Printf("Failed to resolve binary folder to absolute path: %s", ef)
				binaryFolder = ""
			}
		}
	}
	if binaryFolder != "" {
		bindDirs = append(bindDirs, utils.ConvertToDockerBind(binaryFolder)+":"+"/var/lib/pufferpanel/binaries:ro")
	}

	for k, v := range d.Binds {
		bindDirs = append(bindDirs, utils.ConvertToDockerBind(k)+":"+v)
	}

	baseConfig := d.HostConfig

	hostConfig := &baseConfig
	hostConfig.AutoRemove = true
	if hostConfig.NetworkMode == "" {
		hostConfig.NetworkMode = container.NetworkMode(utils.ReplaceTokens(d.Network, data.Variables, utils.PlainReplace))
	}

	hostConfig.Binds = append(hostConfig.Binds, bindDirs...)

	_, hostConfig.PortBindings, err = nat.ParsePortSpecs(utils.ReplaceTokensInArr(d.Ports, data.Variables))
	if err != nil {
		return err
	}

	if hostConfig.PortBindings == nil {
		hostConfig.PortBindings = nat.PortMap{}
	}

	if data.StdInConfig.Port != "" {
		if _, exists := hostConfig.PortBindings[nat.Port(data.StdInConfig.Port+"/tcp")]; !exists {
			//we have a port defined for stdin, we need to also export it
			hostConfig.PortBindings[nat.Port(data.StdInConfig.Port+"/tcp")] = []nat.PortBinding{{
				HostIP: "127.0.0.1", HostPort: data.StdInConfig.Port,
			}}
		}
	}

	if containerConfig.ExposedPorts == nil {
		containerConfig.ExposedPorts = make(nat.PortSet)
	}

	for k := range hostConfig.PortBindings {
		containerConfig.ExposedPorts[k] = struct{}{}
	}

	networkConfig := &network.NetworkingConfig{}

	//for now, default to linux across the board. This resolves problems that Windows has when you use it and docker
	_, err = d.cli.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, &v1.Platform{OS: "linux"}, environment.Server.Id())
	return err
}

func (d *Docker) SendCodeImpl(environment *pufferpanel.Environment, code int) error {
	running, err := environment.IsRunning()

	if err != nil || !running {
		return err
	}

	dockerClient, err := pufferpanel.GetDockerClient()

	if err != nil {
		return err
	}

	ctx := context.Background()
	return dockerClient.ContainerKill(ctx, environment.Server.Id(), cast.ToString(code))
}

func (d *Docker) GetUidImpl(environment *pufferpanel.Environment) int {
	user := d.Config.User
	if user == "" {
		return -1
	}
	return cast.ToInt(strings.Split(user, ":")[0])
}

func (d *Docker) GetGidImpl(environment *pufferpanel.Environment) int {
	user := d.Config.User
	if user == "" {
		return -1
	}
	return cast.ToInt(strings.Split(user, ":")[1])
}

func (d *Docker) handleClose(environment *pufferpanel.Environment, client *client.Client, callback func(int)) {
	exitCode := -1
	okChan, errChan := client.ContainerWait(context.Background(), environment.Server.Id(), container.WaitConditionRemoved)

	select {
	case chanErr := <-errChan:
		{
			exitCode = -999
			environment.Log(logging.Error, "Error from error channel: %s\n", chanErr.Error())
		}
	case info := <-okChan:
		{
			exitCode = cast.ToInt(info.StatusCode)
			if info.Error != nil {
				environment.Log(logging.Error, "Error from info channel: %s\n", info.Error.Message)
			}
		}
	}

	environment.LastExitCode = exitCode
	_ = environment.StatusTracker.WriteMessage(pufferpanel.Transmission{
		Message: pufferpanel.ServerRunning{
			Running:    false,
			Installing: environment.IsInstalling(),
		},
		Type: pufferpanel.MessageTypeStatus,
	})

	_ = environment.Console.Close()
	d.disableSpecialStats = false

	environment.Wait.Unlock()

	if callback != nil {
		callback(exitCode)
	}
}
