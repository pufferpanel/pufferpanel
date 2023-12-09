package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/messages"
	"github.com/spf13/cast"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Docker struct {
	*pufferpanel.BaseEnvironment
	ContainerId   string              `json:"-"`
	ImageName     string              `json:"image"`
	Binds         map[string]string   `json:"bindings,omitempty"`
	Network       string              `json:"networkName,omitempty"`
	Ports         []string            `json:"portBindings,omitempty"`
	Resources     container.Resources `json:"resources,omitempty"`
	Labels        map[string]string   `json:"labels,omitempty"`
	ContainerRoot string              `json:"containerRoot,omitempty"`

	connection       types.HijackedResponse
	cli              *client.Client
	downloadingImage bool
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

	cfg := types.ContainerAttachOptions{
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

	startOpts := types.ContainerStartOptions{}

	msg := messages.Status{Running: true, Installing: d.IsInstalling()}
	_ = d.StatusTracker.WriteMessage(msg)

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
		return &pufferpanel.ServerStats{
			Cpu:    0,
			Memory: 0,
		}, nil
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

	return &pufferpanel.ServerStats{
		Memory: calculateMemoryPercent(data),
		Cpu:    calculateCPUPercent(data),
	}, nil
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
	opts := types.ContainerListOptions{
		Filters: filters.NewArgs(),
	}

	opts.All = true
	opts.Filters.Add("name", d.ContainerId)

	existingContainers, err := client.ContainerList(ctx, opts)

	if len(existingContainers) == 0 {
		return false, err
	} else {
		return true, err
	}
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

		opts := types.ImageListOptions{
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

	op := types.ImagePullOptions{}

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

	//newEnv := os.Environ()
	newEnv := []string{"HOME=" + containerRoot}

	for k, v := range data.Environment {
		newEnv = append(newEnv, fmt.Sprintf("%s=%s", k, v))
	}

	workDir := data.WorkingDirectory
	if workDir == "" {
		workDir = containerRoot
	}

	binaryFolder := config.BinariesFolder.Value()
	if !filepath.IsAbs(binaryFolder) {
		var ef error
		binaryFolder, ef = filepath.Abs(binaryFolder)
		if ef != nil {
			logging.Error.Printf("Failed to resolve binary folder to absolute path: %s", ef)
			binaryFolder = ""
		}
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
		labels[k] = v
	}

	containerConfig := &container.Config{
		AttachStderr:    true,
		AttachStdin:     true,
		AttachStdout:    true,
		Tty:             true,
		OpenStdin:       true,
		NetworkDisabled: false,
		Image:           imageName,
		WorkingDir:      workDir,
		Env:             newEnv,
		Labels:          labels,
	}

	if len(cmdSlice) > 0 {
		containerConfig.Entrypoint = cmdSlice
	}

	if runtime.GOOS != "windows" {
		containerConfig.User = fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid())
	}

	dir := d.GetRootDirectory()

	//convert root dir to a full path, so we can bind it
	if !filepath.IsAbs(dir) {
		dir, err = filepath.Abs(dir)
		if err != nil {
			return err
		}
	}

	bindDirs := []string{convertToBind(dir) + ":" + containerRoot}
	if binaryFolder != "" {
		//bindDirs = append(bindDirs, convertToBind(binaryFolder)+":"+convertToBind(binaryFolder))
		bindDirs = append(bindDirs, convertToBind(binaryFolder)+":"+"/var/lib/pufferpanel/binaries")
	}

	for k, v := range d.Binds {
		bindDirs = append(bindDirs, convertToBind(k)+":"+v)
	}

	hostConfig := &container.HostConfig{
		AutoRemove:   true,
		NetworkMode:  container.NetworkMode(d.Network),
		Resources:    d.Resources,
		Binds:        bindDirs,
		PortBindings: nat.PortMap{},
	}

	networkConfig := &network.NetworkingConfig{}

	_, bindings, err := nat.ParsePortSpecs(d.Ports)
	if err != nil {
		return err
	}
	hostConfig.PortBindings = bindings

	exposedPorts := make(nat.PortSet)
	for k := range bindings {
		exposedPorts[k] = struct{}{}
	}
	containerConfig.ExposedPorts = exposedPorts

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

func (d *Docker) handleClose(client *client.Client, callback func(int)) {
	exitCode := -1
	okChan, errChan := client.ContainerWait(context.Background(), d.ContainerId, container.WaitConditionNotRunning)

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

	msg := messages.Status{Running: false, Installing: d.IsInstalling()}
	_ = d.StatusTracker.WriteMessage(msg)

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
	fullPath = "//" + fullPath
	return fullPath
}
