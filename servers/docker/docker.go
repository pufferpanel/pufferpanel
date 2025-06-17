package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3/utils"
	"slices"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/spf13/cast"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
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

	connection       types.HijackedResponse
	cli              *client.Client
	downloadingImage bool
	statLocker       sync.Mutex
	lastStats        *pufferpanel.ServerStats
	lastStatTime     time.Time
}

func (d *Docker) ExecuteAsyncImpl(environment *pufferpanel.Environment, steps pufferpanel.ExecutionData) error {
	if d.downloadingImage {
		return pufferpanel.ErrImageDownloading
	}

	var err error
	var dockerClient *client.Client
	dockerClient, err = d.getClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	//TODO: This logic may not work anymore, it's complicated to use an existing container with install/uninstall
	exists, err := doesContainerExist(dockerClient, environment.ServerId, ctx)
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

	cfg := container.AttachOptions{
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Stream: true,
	}

	d.connection, err = dockerClient.ContainerAttach(ctx, environment.ServerId, cfg)
	if err != nil {
		return err
	}

	environment.Wait.Add(1)

	go func() {
		defer d.connection.Close()
		_, _ = io.Copy(environment.Wrapper, d.connection.Reader)
	}()

	environment.CreateConsoleStdinProxy(steps.StdInConfig, d.connection.Conn)
	environment.Console.Start()

	go d.handleClose(environment, dockerClient, steps.Callback)

	startOpts := container.StartOptions{}

	_ = environment.StatusTracker.WriteMessage(pufferpanel.Transmission{
		Message: pufferpanel.ServerRunning{
			Running:    true,
			Installing: environment.IsInstalling(),
		},
		Type: pufferpanel.MessageTypeStatus,
	})

	environment.DisplayToConsole(true, "Starting container\n")
	err = dockerClient.ContainerStart(ctx, environment.ServerId, startOpts)
	if err != nil {
		return err
	}

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

	dockerClient, err := d.getClient()
	if err != nil {
		return err
	}
	err = dockerClient.ContainerKill(context.Background(), environment.ServerId, "SIGKILL")
	return err
}

func (d *Docker) IsRunningImpl(environment *pufferpanel.Environment) (bool, error) {
	dockerClient, err := d.getClient()
	if err != nil {
		return false, err
	}

	ctx := context.Background()

	exists, err := doesContainerExist(dockerClient, environment.ServerId, ctx)
	if !exists {
		return false, err
	}

	stats, err := dockerClient.ContainerInspect(ctx, environment.ServerId)
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

		if environment.Server.Stats.Type == "jcmd" {
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

	dockerClient, err := d.getClient()

	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	res, err := dockerClient.ContainerStats(ctx, environment.ServerId, false)
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
		Memory: calculateMemoryPercent(data),
		Cpu:    calculateCPUPercent(data),
	}

	if environment.Server.Stats.Type == "jcmd" {
		cmd, _ := environment.Server.Stats.Metadata["cmd"].(string)
		if cmd == "" {
			cmd = "jcmd"
		}

		r, e := dockerClient.ContainerExecCreate(context.Background(), environment.ServerId, container.ExecOptions{
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

func (d *Docker) getClient() (*client.Client, error) {
	var err error = nil
	if d.cli == nil {
		d.cli, err = client.NewClientWithOpts(client.FromEnv)
		ctx := context.Background()
		d.cli.NegotiateAPIVersion(ctx)
	}
	return d.cli, err
}

func doesContainerExist(client *client.Client, id string, ctx context.Context) (bool, error) {
	opts := container.ListOptions{
		Filters: filters.NewArgs(),
	}

	opts.All = true
	opts.Filters.Add("name", id)

	existingContainers, err := client.ContainerList(ctx, opts)
	if err != nil {
		return false, err
	}

	for _, v := range existingContainers {
		if slices.Contains(v.Names, "/"+id) {
			return true, nil
		}
	}

	return false, nil
}

func (d *Docker) PullImage(environment *pufferpanel.Environment, ctx context.Context, imageName string, force bool) error {
	if d.downloadingImage {
		return pufferpanel.ErrImageDownloading
	}

	if !force {
		exists := false

		parts := strings.SplitN(imageName, ":", 2)
		if len(parts) != 2 {
			imageName = imageName + ":latest"
		}

		opts := image.ListOptions{
			All:     true,
			Filters: filters.NewArgs(),
		}
		opts.Filters.Add("reference", imageName)
		images, err := d.cli.ImageList(ctx, opts)

		if err != nil {
			return err
		}

		for _, v := range images {
			for _, z := range v.RepoTags {
				if z == imageName {
					exists = true
					break
				}
			}
			if exists {
				break
			}
		}

		environment.Log(logging.Debug, "Does image %v exist? %v", imageName, exists)

		if exists {
			return nil
		}
	}

	op := image.PullOptions{}

	environment.Log(logging.Debug, "Downloading image %v", imageName)
	environment.DisplayToConsole(true, "Downloading image for container, please wait\n")

	d.downloadingImage = true
	defer func() {
		d.downloadingImage = false
	}()

	r, err := d.cli.ImagePull(ctx, imageName, op)
	defer utils.Close(r)
	if err != nil {
		return err
	}

	w := &ImageWriter{Parent: environment.ConsoleTracker}
	_, err = io.Copy(w, r)

	if err != nil {
		return err
	}

	environment.Log(logging.Debug, "Downloaded image %v", imageName)
	environment.DisplayToConsole(true, "Downloaded image for container\n")
	return err
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

	imageName := utils.ReplaceTokens(d.ImageName, data.Variables)

	err := d.PullImage(environment, ctx, imageName, false)

	if err != nil {
		return err
	}

	cmdSlice := strslice.StrSlice{}
	if data.Command != "" {
		cmdSlice = append(cmdSlice, data.Command)
	}
	for _, v := range data.Arguments {
		cmdSlice = append(cmdSlice, v)
	}

	environment.Log(logging.Debug, "Container command: %s\n", cmdSlice)

	labels := map[string]string{
		"pufferpanel.server": environment.ServerId,
	}

	for k, v := range d.Labels {
		labels[utils.ReplaceTokens(k, data.Variables)] = utils.ReplaceTokens(v, data.Variables)
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
		containerConfig.Env = append(containerConfig.Env, fmt.Sprintf("%s=%s", k, utils.ReplaceTokens(v, data.Variables)))
	}

	if len(containerConfig.Entrypoint) == 0 && len(cmdSlice) > 0 {
		containerConfig.Entrypoint = cmdSlice
	}

	if containerConfig.User == "" && runtime.GOOS != "windows" {
		containerConfig.User = fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid())
	}

	var dir string
	if containerMountSource != "" {
		dir = filepath.Join(containerMountSource, "servers", environment.ServerId)
	} else {
		dir = environment.GetRootDirectory()
	}

	//convert root dir to a full path, so we can bind it
	if !filepath.IsAbs(dir) {
		dir, err = filepath.Abs(dir)
		if err != nil {
			return err
		}
	}

	bindDirs := []string{convertToBind(dir) + ":" + containerRoot}

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
		bindDirs = append(bindDirs, convertToBind(binaryFolder)+":"+"/var/lib/pufferpanel/binaries")
	}

	for k, v := range d.Binds {
		bindDirs = append(bindDirs, convertToBind(k)+":"+v)
	}

	baseConfig := d.HostConfig

	hostConfig := &baseConfig
	hostConfig.AutoRemove = true
	if hostConfig.NetworkMode == "" {
		hostConfig.NetworkMode = container.NetworkMode(utils.ReplaceTokens(d.Network, data.Variables))
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
	_, err = d.cli.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, &v1.Platform{OS: "linux"}, environment.ServerId)
	return err
}

func (d *Docker) SendCodeImpl(environment *pufferpanel.Environment, code int) error {
	running, err := environment.IsRunning()

	if err != nil || !running {
		return err
	}

	dockerClient, err := d.getClient()

	if err != nil {
		return err
	}

	ctx := context.Background()
	return dockerClient.ContainerKill(ctx, environment.ServerId, cast.ToString(code))
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
	okChan, errChan := client.ContainerWait(context.Background(), environment.ServerId, container.WaitConditionRemoved)

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

	environment.Wait.Done()

	_ = environment.StatusTracker.WriteMessage(pufferpanel.Transmission{
		Message: pufferpanel.ServerRunning{
			Running:    false,
			Installing: environment.IsInstalling(),
		},
		Type: pufferpanel.MessageTypeStatus,
	})

	_ = environment.Console.Close()

	if callback != nil {
		callback(exitCode)
	}
}

func calculateCPUPercent(v *container.StatsResponse) float64 {
	//this math is from https://docs.docker.com/reference/api/engine/version/v1.45/#tag/Container/operation/ContainerStats
	cpuDelta := v.CPUStats.CPUUsage.TotalUsage - v.PreCPUStats.CPUUsage.TotalUsage
	systemCpuDelta := v.CPUStats.SystemUsage - v.PreCPUStats.SystemUsage
	numCpus := int(v.CPUStats.OnlineCPUs)
	if numCpus == 0 {
		numCpus = len(v.CPUStats.CPUUsage.PercpuUsage)
	}
	return (float64(cpuDelta) / float64(systemCpuDelta)) * float64(numCpus) * 100.0
}

func calculateMemoryPercent(v *container.StatsResponse) float64 {
	return float64(v.MemoryStats.Usage)
}

func convertToBind(source string) string {
	fullPath, err := filepath.Abs(source)
	if err != nil {
		panic(err)
	}

	fullPath = strings.ReplaceAll(fullPath, "\\", "/")
	fullPath = strings.ReplaceAll(fullPath, ":", "")
	//lowercase first character as that's the drive
	fullPath = strings.ToLower(string(fullPath[0])) + fullPath[1:]
	fullPath = "/" + fullPath
	return fullPath
}
