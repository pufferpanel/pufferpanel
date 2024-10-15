package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	*pufferpanel.BaseEnvironment
	ContainerId   string               `json:"-"`
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

func (d *Docker) dockerExecuteAsync(steps pufferpanel.ExecutionData) error {
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
	exists, err := d.doesContainerExist(dockerClient, ctx)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("docker container already exists")
	}

	err = d.createContainer(ctx, steps)
	if err != nil {
		return err
	}

	cfg := container.AttachOptions{
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Stream: true,
	}

	d.connection, err = dockerClient.ContainerAttach(ctx, d.ContainerId, cfg)
	if err != nil {
		return err
	}

	d.Wait.Add(1)

	go func() {
		defer d.connection.Close()
		_, _ = io.Copy(d.Wrapper, d.connection.Reader)
	}()

	d.BaseEnvironment.CreateConsoleStdinProxy(steps.StdInConfig, d.connection.Conn)
	d.BaseEnvironment.Console.Start()

	go d.handleClose(dockerClient, steps.Callback)

	startOpts := container.StartOptions{}

	_ = d.StatusTracker.WriteMessage(pufferpanel.Transmission{
		Message: pufferpanel.ServerRunning{
			Running:    true,
			Installing: d.IsInstalling(),
		},
		Type: pufferpanel.MessageTypeStatus,
	})

	d.DisplayToConsole(true, "Starting container\n")
	err = dockerClient.ContainerStart(ctx, d.ContainerId, startOpts)
	if err != nil {
		return err
	}

	return err
}

func (d *Docker) kill() (err error) {
	running, err := d.IsRunning()
	if err != nil {
		return err
	}

	if !running {
		return
	}

	dockerClient, err := d.getClient()
	if err != nil {
		return err
	}
	err = dockerClient.ContainerKill(context.Background(), d.ContainerId, "SIGKILL")
	return
}

func (d *Docker) isRunning() (bool, error) {
	dockerClient, err := d.getClient()
	if err != nil {
		return false, err
	}

	ctx := context.Background()

	exists, err := d.doesContainerExist(dockerClient, ctx)
	if !exists {
		return false, err
	}

	stats, err := dockerClient.ContainerInspect(ctx, d.ContainerId)
	if err != nil {
		return false, err
	}
	return stats.State.Running, nil
}

func (d *Docker) GetStats() (*pufferpanel.ServerStats, error) {
	running, err := d.IsRunning()
	if err != nil {
		return nil, err
	}

	if !running {
		stats := &pufferpanel.ServerStats{
			Cpu:    0,
			Memory: 0,
		}

		if d.Server.Stats.Type == "jcmd" {
			stats.Jvm = &pufferpanel.JvmStats{}
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
	res, err := dockerClient.ContainerStats(ctx, d.ContainerId, false)
	defer func() {
		if res.Body != nil {
			pufferpanel.Close(res.Body)
		}
	}()
	if err != nil {
		return nil, err
	}

	data := &types.StatsJSON{}
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

	if d.Server.Stats.Type == "jcmd" {
		cmd, _ := d.Server.Stats.Metadata["cmd"].(string)
		if cmd == "" {
			cmd = "jcmd"
		}

		r, e := dockerClient.ContainerExecCreate(context.Background(), d.ContainerId, types.ExecConfig{
			AttachStderr: true,
			AttachStdout: true,
			Cmd:          []string{cmd, "1", "GC.heap_info"},
		})

		if e == nil {
			rw, e := dockerClient.ContainerExecAttach(context.Background(), r.ID, types.ExecStartCheck{
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

				stats.Jvm = pufferpanel.ParseJCMDResponse(jcmdData)
			}
		}
		if stats.Jvm == nil {
			stats.Jvm = &pufferpanel.JvmStats{}
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

func (d *Docker) doesContainerExist(client *client.Client, ctx context.Context) (bool, error) {
	opts := container.ListOptions{
		Filters: filters.NewArgs(),
	}

	opts.All = true
	opts.Filters.Add("name", d.ContainerId)

	existingContainers, err := client.ContainerList(ctx, opts)
	if err != nil {
		return false, err
	}

	for _, v := range existingContainers {
		if slices.Contains(v.Names, "/"+d.ContainerId) {
			return true, nil
		}
	}

	return false, nil
}

func (d *Docker) PullImage(ctx context.Context, imageName string, force bool) error {
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

		d.Log(logging.Debug, "Does image %v exist? %v", imageName, exists)

		if exists {
			return nil
		}
	}

	op := image.PullOptions{}

	d.Log(logging.Debug, "Downloading image %v", imageName)
	d.DisplayToConsole(true, "Downloading image for container, please wait\n")

	d.downloadingImage = true
	defer func() {
		d.downloadingImage = false
	}()

	r, err := d.cli.ImagePull(ctx, imageName, op)
	defer pufferpanel.Close(r)
	if err != nil {
		return err
	}

	w := &ImageWriter{Parent: d.ConsoleTracker}
	_, err = io.Copy(w, r)

	if err != nil {
		return err
	}

	d.Log(logging.Debug, "Downloaded image %v", imageName)
	d.DisplayToConsole(true, "Downloaded image for container\n")
	return err
}

func (d *Docker) createContainer(ctx context.Context, data pufferpanel.ExecutionData) error {
	d.Log(logging.Debug, "Creating container")
	containerRoot := d.ContainerRoot
	if containerRoot == "" {
		containerRoot = "/pufferpanel"
	}

	if runtime.GOOS != "windows" {
		if !filepath.IsAbs(containerRoot) {
			return pufferpanel.ErrPathNotAbs(containerRoot)
		}
	}

	imageName := pufferpanel.ReplaceTokens(d.ImageName, data.Variables)

	err := d.PullImage(ctx, imageName, false)

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

	d.Log(logging.Debug, "Container command: %s\n", cmdSlice)

	labels := map[string]string{
		"pufferpanel.server": d.ContainerId,
	}

	for k, v := range d.Labels {
		labels[pufferpanel.ReplaceTokens(k, data.Variables)] = pufferpanel.ReplaceTokens(v, data.Variables)
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
		containerConfig.Env = append(containerConfig.Env, fmt.Sprintf("%s=%s", k, v))
	}

	if len(containerConfig.Entrypoint) == 0 && len(cmdSlice) > 0 {
		containerConfig.Entrypoint = cmdSlice
	}

	if containerConfig.User == "" && runtime.GOOS != "windows" {
		containerConfig.User = fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid())
	}

	var dir string
	if config.DockerRootPath.Value() != "" {
		dir = filepath.Join(config.DockerRootPath.Value(), "servers", d.ServerId)
	} else {
		dir = d.GetRootDirectory()
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
	if config.DockerRootPath.Value() != "" {
		binaryFolder = filepath.Join(config.DockerRootPath.Value(), "binaries")
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
		hostConfig.NetworkMode = container.NetworkMode(d.Network)
	}

	hostConfig.Binds = append(hostConfig.Binds, bindDirs...)

	if hostConfig.PortBindings == nil {
		hostConfig.PortBindings = nat.PortMap{}
	}

	_, hostConfig.PortBindings, err = nat.ParsePortSpecs(d.Ports)
	if err != nil {
		return err
	}

	if containerConfig.ExposedPorts == nil {
		containerConfig.ExposedPorts = make(nat.PortSet)
	}

	for k := range hostConfig.PortBindings {
		containerConfig.ExposedPorts[k] = struct{}{}
	}

	networkConfig := &network.NetworkingConfig{}

	//for now, default to linux across the board. This resolves problems that Windows has when you use it and docker
	_, err = d.cli.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, &v1.Platform{OS: "linux"}, d.ContainerId)
	return err
}

func (d *Docker) SendCode(code int) error {
	running, err := d.IsRunning()

	if err != nil || !running {
		return err
	}

	dockerClient, err := d.getClient()

	if err != nil {
		return err
	}

	ctx := context.Background()
	return dockerClient.ContainerKill(ctx, d.ContainerId, cast.ToString(code))
}

func (d *Docker) GetUid() int {
	user := d.Config.User
	if user == "" {
		return -1
	}
	return cast.ToInt(strings.Split(user, ":")[0])
}

func (d *Docker) GetGid() int {
	user := d.Config.User
	if user == "" {
		return -1
	}
	return cast.ToInt(strings.Split(user, ":")[1])
}

func (d *Docker) handleClose(client *client.Client, callback func(int)) {
	exitCode := -1
	okChan, errChan := client.ContainerWait(context.Background(), d.ContainerId, container.WaitConditionRemoved)

	select {
	case chanErr := <-errChan:
		{
			exitCode = -999
			d.Log(logging.Error, "Error from error channel: %s\n", chanErr.Error())
		}
	case info := <-okChan:
		{
			exitCode = cast.ToInt(info.StatusCode)
			if info.Error != nil {
				d.Log(logging.Error, "Error from info channel: %s\n", info.Error.Message)
			}
		}
	}

	d.LastExitCode = exitCode

	d.Wait.Done()

	_ = d.StatusTracker.WriteMessage(pufferpanel.Transmission{
		Message: pufferpanel.ServerRunning{
			Running:    false,
			Installing: d.IsInstalling(),
		},
		Type: pufferpanel.MessageTypeStatus,
	})

	_ = d.BaseEnvironment.Console.Close()

	if callback != nil {
		callback(exitCode)
	}
}

func calculateCPUPercent(v *types.StatsJSON) float64 {
	// Max number of 100ns intervals between the previous time read and now
	possIntervals := uint64(v.Read.Sub(v.PreRead).Nanoseconds()) // Start with number of ns intervals
	possIntervals /= 100                                         // Convert to number of 100ns intervals
	//possIntervals *= uint64(v.NumProcs)                          // Multiple by the number of processors

	// Intervals used
	intervalsUsed := v.CPUStats.CPUUsage.TotalUsage - v.PreCPUStats.CPUUsage.TotalUsage

	// Percentage avoiding divide-by-zero
	if possIntervals > 0 {
		return float64(intervalsUsed) / float64(possIntervals)
	}
	return 0.00
}

func calculateMemoryPercent(v *types.StatsJSON) float64 {
	return float64(v.MemoryStats.Usage)
}

func convertToBind(source string) string {
	if runtime.GOOS != "windows" {
		return source
	}
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
