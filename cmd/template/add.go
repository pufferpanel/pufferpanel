/*
 Copyright 2020 Padduck, LLC
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

package template

import (
	"fmt"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"github.com/spf13/cobra"
	"os"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add template from disk",
	Run:   runAdd,
	Args:  cobra.ExactArgs(1),
}

var templateName string
var readme string

func init() {
	AddCmd.Flags().StringVar(&templateName, "name", "", "process id of daemon")
	AddCmd.Flags().StringVar(&readme, "readme", "", "path to readme file")
	_ = AddCmd.MarkFlagFilename("name", "*.json")
	_ = AddCmd.MarkFlagFilename("readme", "*.md")
}

func runAdd(cmd *cobra.Command, args []string) {
	err := pufferpanel.LoadConfig("")
	if err != nil {
		fmt.Printf("Error loading config: %s\n", err.Error())
		os.Exit(1)
		return
	}

	db, err := database.GetConnection()
	if err != nil {
		fmt.Printf("Error connecting to database: %s\n", err.Error())
		return
	}

	ts := &services.Template{DB: db}

	err = ts.ImportTemplate(templateName, args[0], readme)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
		return
	}
}
