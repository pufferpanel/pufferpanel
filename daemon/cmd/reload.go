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
)

var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reloads pufferd",
	Run: func(cmd *cobra.Command, args []string) {
		err := runReload()
		if err != nil {
			logging.Exception("error running reload", err)
		}
	},
}

var reloadPid int

func init() {
	reloadCmd.Flags().IntVar(&reloadPid, "pid", 0, "process id of daemon")
	reloadCmd.MarkPersistentFlagRequired("pid")
}

func runReload() error {
	proc, err := os.FindProcess(reloadPid)
	if err != nil || proc == nil {
		if err == nil && proc == nil {
			err = errors.New("no process found")
		}
		return err
	}
	return proc.Signal(syscall.Signal(1))
}
