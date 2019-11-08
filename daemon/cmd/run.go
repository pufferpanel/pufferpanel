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
	"fmt"
	"github.com/pufferpanel/pufferpanel/v2/daemon"
	"github.com/pufferpanel/pufferpanel/v2/daemon/entry"
	"github.com/pufferpanel/pufferpanel/v2/shared/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the daemon",
	Run:   runRun,
}

func runRun(cmd *cobra.Command, args []string) {
	daemon.SetDefaults()
	_ = daemon.LoadConfig()

	var logPath = viper.GetString("data.logs")
	_ = logging.WithLogDirectory(logPath, logging.DEBUG, nil)

	err := <-entry.Start()
	if err != nil {
		fmt.Printf("Error running: %s", err.Error())
	}
}
