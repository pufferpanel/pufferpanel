/*
 Copyright 2019 Padduck, LLC

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

package main

import (
	"errors"
	"github.com/pufferpanel/pufferpanel/v2/shared/logging"
	"github.com/spf13/cobra"
	"os"
	"syscall"
	"time"
)

var shutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "Shuts down pufferd",
	Run: func(cmd *cobra.Command, args []string) {
		err := runShutdown()
		if err != nil {
			logging.Exception("error running shutdown", err)
		}
	},
}

var shutdownPid int

func init() {
	shutdownCmd.Flags().IntVar(&shutdownPid, "pid", 0, "process id of daemon")
	shutdownCmd.MarkFlagRequired("pid")
}

func runShutdown() error {
	proc, err := os.FindProcess(shutdownPid)
	if err != nil || proc == nil {
		if err == nil && proc == nil {
			err = errors.New("no process found")
		}
		return err
	}
	err = proc.Signal(syscall.Signal(15))
	if err != nil {
		return err
	}

	wait := make(chan error)

	waitForProcess(proc, wait)

	err = <-wait

	if err != nil {
		return err
	}

	err = proc.Release()
	if err != nil {
		return err
	}

	return nil
}

func waitForProcess(process *os.Process, c chan error) {
	var err error
	timer := time.NewTicker(100 * time.Millisecond)
	go func() {
		for range timer.C {
			err = process.Signal(syscall.Signal(0))
			if err != nil {
				if err.Error() == "os: process already finished" {
					c <- nil
				} else {
					c <- err
				}

				timer.Stop()
			} else {
			}
		}
	}()
}
