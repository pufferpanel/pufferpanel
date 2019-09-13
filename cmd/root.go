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
	"github.com/pufferpanel/apufferi/v3/logging"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "pufferpanel",
	Short: "Game Server Management Panel",
	Run:   root,
}

func init() {
	rootCmd.AddCommand(
		runCmd,
		versionCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func root(cmd *cobra.Command, args []string) {
	err := executeRun(cmd, args)
	if err != nil {
		logging.Exception("An error has occurred while executing", err)
	}
}
