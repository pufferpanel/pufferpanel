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

package lxd

import (
	"github.com/docker/docker/pkg/ioutils"
	client "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"github.com/pufferpanel/pufferpanel/v2/shared"
	"github.com/pufferpanel/pufferpanel/v2/shared/logging"
	"github.com/pufferpanel/pufferpanel/v2/daemon"
	"github.com/pufferpanel/pufferpanel/v2/daemon/environments/envs"
	"io"
	"os"
	"sync"
	"time"
)

type lxc struct {
	*envs.BaseEnvironment
	shared.TypeWithMetadata
	ContainerId string `json:"-"`
	ImageName   string `json:"image"`

	connection   client.InstanceServer
	wait         *sync.WaitGroup
	stdin        io.Writer
	stdinReader  io.ReadCloser
	stdout       io.Reader
	stdoutWriter io.WriteCloser
}

func (l *lxc) executeAsync(cmd string, args []string, env map[string]string, callback func(graceful bool)) error {
	running, err := l.IsRunning()
	if err != nil {
		return err
	}
	if running {
		return daemon.ErrContainerRunning
	}

	//send an exec
	req := api.ContainerExecPost{
		Command:     append([]string{cmd}, args...),
		WaitForWS:   true,
		Interactive: true,
	}

	l.stdinReader, l.stdin = io.Pipe()
	l.stdout, l.stdoutWriter = io.Pipe()

	out := l.CreateWrapper()
	l.stdoutWriter = ioutils.NopWriteCloser(io.MultiWriter(out, l.stdoutWriter))

	finished := make(chan bool)

	execArgs := &client.ContainerExecArgs{
		Stdin:    l.stdinReader,
		Stdout:   l.stdoutWriter,
		Stderr:   l.stdoutWriter,
		DataDone: finished,
	}

	l.wait.Add(1)

	_, _, err = l.connection.GetContainer(l.ContainerId)
	if err != nil {
		logging.Devel("Container does not exist (%s)", err.Error())
		//container no exist, make it
		req := api.ContainersPost{
			Name: l.ContainerId,
			Source: api.ContainerSource{
				Type:     "image",
				Alias:    l.ImageName,
				Server:   "https://images.linuxcontainers.org",
				Protocol: "simplestreams",
			},
			ContainerPut: api.ContainerPut{
				Ephemeral: true,
			},
		}

		// Get LXD to create the container (background operation)
		op, err := l.connection.CreateContainer(req)
		if err != nil {
			return err
		}

		// Wait for the operation to complete
		err = op.Wait()
		if err != nil {
			return err
		}
	}

	// Get LXD to start the container (background operation)
	reqState := api.ContainerStatePut{
		Action:  "start",
		Timeout: -1,
	}

	op, err := l.connection.UpdateContainerState(l.ContainerId, reqState, "")
	if err != nil {
		return err
	}

	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		return err
	}

	op, err = l.connection.ExecContainer(l.ContainerId, req, execArgs)
	if err != nil {
		return err
	}

	go func() {
		defer l.wait.Done()
		<-finished
		var internalError error
		defer func() {
			if internalError != nil {
				logging.Exception("Error stopping container "+l.ContainerId, internalError)
			}
			if callback != nil {
				callback(internalError == nil)
			}
		}()

		/*internalError = op.Wait()
		if internalError != nil {
			return
		}*/

		op2, _ := l.connection.UpdateContainerState(l.ContainerId, api.ContainerStatePut{
			Action:  "stop",
			Timeout: -1,
			Force:   true,
		}, "")

		if op2 != nil {
			internalError = op2.Wait()
		}
	}()

	return nil
}

func (l *lxc) ExecuteInMainProcess(cmd string) (err error) {
	running, err := l.IsRunning()
	if err != nil {
		return
	}
	if !running {
		err = daemon.ErrServerOffline
		return
	}

	_, _ = l.stdin.Write([]byte(cmd + "\n"))
	return
}

func (l *lxc) IsRunning() (bool, error) {
	logging.Debug("Checking if %s is running", l.ContainerId)
	c, err := l.getConnection()
	if err != nil {
		return false, err
	}

	state, _, err := c.GetContainerState(l.ContainerId)
	//since errors can mean it doesn't exist, and they don't expose a way to tell, we have to ignore all errors
	if err != nil {
		return false, nil
	}

	if state.StatusCode == api.Running {
		return true, nil
	}

	return false, nil
}

func (l *lxc) Create() error {
	err := os.Mkdir(l.RootDirectory, 0755)
	if err != nil {
		return err
	}

	return err
}

func (l *lxc) GetStats() (*daemon.ServerStats, error) {
	return &daemon.ServerStats{
		Cpu:    0,
		Memory: 0,
	}, nil
}

func (l *lxc) WaitForMainProcess() error {
	return l.WaitForMainProcessFor(0)
}

func (l *lxc) WaitForMainProcessFor(timeout int) (err error) {
	running, err := l.IsRunning()
	if err != nil {
		return
	}
	if running {
		if timeout > 0 {
			var timer = time.AfterFunc(time.Duration(timeout)*time.Millisecond, func() {
				err = l.Kill()
			})
			l.wait.Wait()
			timer.Stop()
		} else {
			l.wait.Wait()
		}
	}
	return
}

func (l *lxc) Kill() error {
	op2, err := l.connection.UpdateContainerState(l.ContainerId, api.ContainerStatePut{
		Action:  "stop",
		Timeout: -1,
		Force:   true,
	}, "")

	if err != nil {
		return err
	}

	return op2.Wait()
}

func (l *lxc) getConnection() (client.InstanceServer, error) {
	var err error
	if l.connection == nil {
		l.connection, err = client.ConnectLXDUnix("", nil)
	}
	return l.connection, err
}
