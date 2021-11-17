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
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/spf13/cobra"
	"os"
	"runtime/debug"
)

func main() {
	if !pufferpanel.UserInGroup("pufferpanel") {
		fmt.Println("You do not have permission to use this command")
		return
	}

	defer logging.Close()

	defer func() {
		if err := recover(); err != nil {
			stacktrace := debug.Stack()
			logging.Error.Printf("%s\n%s\n", err, stacktrace)

			os.Exit(2)
		}
	}()

	Execute()
}

var rootCmd = &cobra.Command{
	Use:   "pufferpanel",
	Short: "Game Server Management Panel",
}

var workDir string

func init() {
	rootCmd.PersistentFlags().StringVar(&workDir, "workDir", "", "Set working directory")

	cobra.OnInitialize(setWorkDir)

	rootCmd.AddCommand(
		runCmd,
		versionCmd,
		userCmd,
		shutdownCmd,
	)
}

func setWorkDir() {
	if workDir != "" {
		err := os.Chdir(workDir)
		if err != nil {
			panic(err)
		}
	}
}

func Execute() {
	rootCmd.SetVersionTemplate(pufferpanel.Display)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
