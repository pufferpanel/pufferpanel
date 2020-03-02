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
	"encoding/json"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

	err = importTemplate(templateName, args[0], readme, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
		return
	}
}

func importTemplate(name, templatePath, readmePath string, service *services.Template) error {
	template, err := openTemplate(templatePath)

	if err != nil {
		return err
	}

	if name == "" {
		name = strings.TrimSuffix(filepath.Base(templatePath), filepath.Ext(templatePath))
	}

	model := &models.Template{
		Template: template,
		Name:     name,
		Readme:   "",
	}

	if readmePath != "" {
		data, err := openReadme(readmePath)
		if err != nil {
			fmt.Printf("No readme located at %s, will still import template\n", readmePath)
			//return err
		}
		model.Readme = data
	}

	if service == nil {
		db, err := database.GetConnection()
		if err != nil {
			return err
		}
		defer pufferpanel.Close(db)

		service = &services.Template{DB: db}
	}

	err = service.Save(model)
	return err
}

func openTemplate(path string) (t pufferpanel.Template, err error) {
	file, err := os.Open(path)
	defer pufferpanel.Close(file)
	if err != nil {
		return
	}

	err = json.NewDecoder(file).Decode(&t)
	return
}

func openReadme(path string) (string, error) {
	file, err := os.Open(path)
	defer pufferpanel.Close(file)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(data), err
}
