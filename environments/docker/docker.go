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
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/config"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/messages"
	"github.com/spf13/cast"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type docker struct {
	*pufferpanel.BaseEnvironment
	ContainerId string              `json:"-"`
	ImageName   string              `json:"image"`
	Binds       map[string]string   `json:"bindings,omitempty"`
	NetworkMode string              `json:"networkMode,omitempty"`
	Network     string              `json:"networkName,omitempty"`
	Ports       []string            `json:"portBindings,omitempty"`
	Resources   container.Resources `json:"resources,omitempty"`

	connection       types.HijackedResponse
	cli              *client.Client
	downloadingImage bool
}

func (d *docker) dockerExecuteAsync(steps pufferpanel.ExecutionData) error {
	running, err := d.IsRunning()
	if err != nil {
		return err
	}
	if running {
		return pufferpanel.ErrContainerRunning
	}

	if d.downloadingImage {
		return pufferpanel.ErrImageDownloading
	}

	dockerClient, err := d.getClient()
	ctx := context.Background()

	//TODO: This logic may not work anymore, it's complicated to use an existing container with install/uninstall
	exists, err := d.doesContainerExist(dockerClient, ctx)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("docker container already exists")
	}

	err = d.createContainer(dockerClient, ctx, steps.Command, steps.Arguments, steps.Environment, steps.WorkingDirectory)
	if err != nil {
		return err
	}

	config := types.ContainerAttachOptions{
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Stream: true,
	}

	d.connection, err = dockerClient.ContainerAttach(ctx, d.ContainerId, config)
	if err != nil {
		return err
	}

	d.Wait.Add(1)

	go func() {
		defer d.connection.Close()
		wrapper := d.CreateWrapper()
		_, _ = io.Copy(wrapper, d.connection.Reader)
		//because we use the auto-delete, we don't manually stop the container
		//c, _ := d.getClient()
		//err = c.ContainerStop(context.Background(), d.ContainerId, nil)
		//if err != nil {
		//	logging.Error.Printf("Error stopping container "+d.ContainerId, err)
		//}

		okChan, errChan := dockerClient.ContainerWait(ctx, d.ContainerId, container.WaitConditionRemoved)
		select {
		case _ = <-okChan:
		case chanErr := <-errChan:
			if chanErr != nil {
				d.Log(logging.Error, "Error from error channel, awaiting exit `%v`\n", chanErr)
			}
		}

		d.Wait.Done()

		msg := messages.Status{Running: false}
		_ = d.WSManager.WriteMessage(msg)

		if steps.Callback != nil {
			steps.Callback(err == nil)
		}
	}()

	startOpts := types.ContainerStartOptions{}

	msg := messages.Status{Running: true}
	_ = d.WSManager.WriteMessage(msg)

	d.DisplayToConsole(true, "Starting container\n")
	err = dockerClient.ContainerStart(ctx, d.ContainerId, startOpts)
	if err != nil {
		return err
	}
	return err
}

func (d *docker) ExecuteInMainProcess(cmd string) (err error) {
	running, err := d.IsRunning()
	if err != nil {
		return
	}
	if !running {
		err = pufferpanel.ErrServerOffline
		return
	}

	_, _ = d.connection.Conn.Write([]byte(cmd + "\n"))
	return
}

func (d *docker) Kill() (err error) {
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

func (d *docker) Create() error {
	/*go func() {
		c, err := d.getClient()
		if err != nil {
			logging.Error.Printf("Error getting docker client: %s\n", err.Error())
			d.DisplayToConsole(true, "Error downloading image")
			return
		}
		err = d.pullImage(c, context.Background(), false)
		if err != nil {
			logging.Error.Printf("Error downloading image: %s\n", err.Error())
			d.DisplayToConsole(true, "Error downloading image")
			return
		}
	}()*/
	return os.Mkdir(d.RootDirectory, 0755)
}

func (d *docker) IsRunning() (bool, error) {
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

func (d *docker) GetStats() (*pufferpanel.ServerStats, error) {
	running, err := d.IsRunning()
	if err != nil {
		return nil, err
	}

	if !running {
		return nil, pufferpanel.ErrServerOffline
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

func (d *docker) WaitForMainProcess() error {
	return d.WaitForMainProcessFor(0)
}

func (d *docker) WaitForMainProcessFor(timeout time.Duration) (err error) {
	running, err := d.IsRunning()
	if err != nil {
		return
	}
	if running {
		if timeout > 0 {
			var timer = time.AfterFunc(timeout, func() {
				err = d.Kill()
			})
			d.Wait.Wait()
			timer.Stop()
		} else {
			d.Wait.Wait()
		}
	}
	return
}

func (d *docker) getClient() (*client.Client, error) {
	var err error = nil
	if d.cli == nil {
		d.cli, err = client.NewClientWithOpts(client.FromEnv)
		ctx := context.Background()
		d.cli.NegotiateAPIVersion(ctx)
	}
	return d.cli, err
}

func (d *docker) doesContainerExist(client *client.Client, ctx context.Context) (bool, error) {
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

func (d *docker) pullImage(client *client.Client, ctx context.Context, force bool) error {
	if d.downloadingImage {
		return pufferpanel.ErrImageDownloading
	}

	exists := false

	opts := types.ImageListOptions{
		All:     true,
		Filters: filters.NewArgs(),
	}
	opts.Filters.Add("reference", d.ImageName)
	images, err := client.ImageList(ctx, opts)

	if err != nil {
		return err
	}

	if len(images) >= 1 {
		exists = true
	}

	d.Log(logging.Debug, "Does image %v exist? %v", d.ImageName, exists)

	if exists && !force {
		return nil
	}

	op := types.ImagePullOptions{}

	d.Log(logging.Debug, "Downloading image %v", d.ImageName)
	d.DisplayToConsole(true, "Downloading image for container, please wait\n")

	d.downloadingImage = true
	defer func() {
		d.downloadingImage = false
	}()

	r, err := client.ImagePull(ctx, d.ImageName, op)
	defer pufferpanel.Close(r)
	if err != nil {
		return err
	}

	w := &ImageWriter{Parent: d.WSManager}
	_, err = io.Copy(w, r)

	if err != nil {
		return err
	}

	d.Log(logging.Debug, "Downloaded image %v", d.ImageName)
	d.DisplayToConsole(true, "Downloaded image for container\n")
	return err
}

func (d *docker) createContainer(client *client.Client, ctx context.Context, cmd string, args []string, env map[string]string, workDir string) error {
	d.Log(logging.Debug, "Creating container")
	containerRoot := "/pufferpanel"
	err := d.pullImage(client, ctx, false)

	if err != nil {
		return err
	}

	//newEnv := os.Environ()
	newEnv := []string{"HOME=" + containerRoot}

	for k, v := range env {
		newEnv = append(newEnv, fmt.Sprintf("%s=%s", k, v))
	}

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
	cmdSlice = append(cmdSlice, cmd)
	for _, v := range args {
		cmdSlice = append(cmdSlice, v)
	}

	d.Log(logging.Debug, "Container command: %s\n", cmdSlice)

	containerConfig := &container.Config{
		AttachStderr:    true,
		AttachStdin:     true,
		AttachStdout:    true,
		Tty:             true,
		OpenStdin:       true,
		NetworkDisabled: false,
		Image:           d.ImageName,
		WorkingDir:      workDir,
		Env:             newEnv,
		Entrypoint:      cmdSlice,
		Labels: map[string]string{
			"pufferpanel.server": d.ContainerId,
		},
	}

	if runtime.GOOS == "linux" {
		containerConfig.User = fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid())
	}

	dir := d.RootDirectory

	//convert root dir to a full path, so we can bind it
	if !filepath.IsAbs(dir) {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		dir = filepath.Join(pwd, dir)
	}

	bindDirs := []string{dir + ":" + containerRoot}
	if binaryFolder != "" {
		bindDirs = append(bindDirs, binaryFolder+":"+binaryFolder)
	}

	hostConfig := &container.HostConfig{
		AutoRemove:   true,
		NetworkMode:  container.NetworkMode(d.NetworkMode),
		Resources:    d.Resources,
		Binds:        bindDirs,
		PortBindings: nat.PortMap{},
	}

	for k, v := range d.Binds {
		hostConfig.Binds = append(hostConfig.Binds, k+":"+v)
	}

	networkConfig := &network.NetworkingConfig{}

	_, bindings, err := nat.ParsePortSpecs(d.Ports)
	if err != nil {
		return err
	}
	hostConfig.PortBindings = bindings

	exposedPorts := make(nat.PortSet)
	for k, _ := range bindings {
		exposedPorts[k] = struct{}{}
	}
	containerConfig.ExposedPorts = exposedPorts

	//for now, default to linux across the board. This resolves problems that Windows has when you use it and docker
	_, err = client.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, &v1.Platform{OS: "linux"}, d.ContainerId)
	return err
}

func (d *docker) SendCode(code int) error {
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
