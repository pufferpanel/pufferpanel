package main

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/legacy"
	"github.com/pufferpanel/pufferpanel/v2/programs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/tools/go/ssa/interp/testdata/src/runtime"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"github.com/otiai10/copy"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrates from v1 to v2",
	Run:   migrate,
}

var migrateConfig string

func init() {

	var defaultPath = "/etc/pufferd/config.json"
	if runtime.GOOS == "windows" {
		defaultPath = "config.json"
	}

	migrateCmd.Flags().StringVarP(&migrateConfig, "config", "c", defaultPath, "Location of old pufferd config")
}

func migrate(cmd *cobra.Command, args []string) {
	confirm := false
	_ = survey.AskOne(&survey.Confirm{
		Message: "Are you SURE you wish to migrate from v1 to v2? There is NO WARRANTY OR GUARANTEE this option will fully migrate your servers.",
	}, &confirm)

	if !confirm {
		return
	}

	oldConfig := &legacy.Config{}
	data, err := ioutil.ReadFile(migrateConfig)
	if err != nil {
		fmt.Printf("Error loading legacy config: %s\n", err)
		os.Exit(1)
		return
	}
	err = json.Unmarshal(data, oldConfig)
	if err != nil {
		fmt.Printf("Error loading legacy config: %s\n", err)
		os.Exit(1)
		return
	}

	err = pufferpanel.LoadConfig("")
	if err != nil {
		fmt.Printf("Error loading new config: %s\n", err)
		os.Exit(1)
		return
	}

	//start migration of data.... begin the hell
	serversFolder := oldConfig.ServerFolder
	if serversFolder == "" {
		serversFolder = "/var/lib/pufferd/servers"
	}

	files, err := ioutil.ReadDir(serversFolder)
	if err != nil {
		fmt.Printf("Error reading directory: %s\n", err)
		os.Exit(1)
		return
	}

	newFolders := viper.GetString("daemon.data.servers")

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		serverId := strings.TrimSuffix(file.Name(), ".json")

		data, err := ioutil.ReadFile(filepath.Join(serversFolder, file.Name()))
		if err != nil {
			fmt.Printf("Error reading server file %s: %s\n", file.Name(), err)
			os.Exit(1)
			return
		}

		legacyData := &legacy.ServerJson{}
		err = json.Unmarshal(data, legacyData)
		if err != nil {
			fmt.Printf("Error reading server file %s: %s\n", file.Name(), err)
			os.Exit(1)
			return
		}

		newServer := &programs.Program{
			Server:             pufferpanel.Server{
				Variables:      map[string]pufferpanel.Variable{},
				Display:        legacyData.ProgramData.Display,
				Environment:    nil,
				Installation:   nil,
				Uninstallation: nil,
				Type:           pufferpanel.Type{Type: legacyData.ProgramData.Type},
				Identifier:     serverId,
				Execution:      pufferpanel.Execution{},
			},
		}

		for k, v := range legacyData.ProgramData.Data {
			newServer.Variables[k] = pufferpanel.Variable{
				Description:  v.Description,
				Display:      v.Display,
				Internal:     v.Internal,
				Required:     v.Required,
				Value:        v.Value,
				UserEditable: v.UserEditable,
				Type:         pufferpanel.Type{Type: "string"},
				Options:      nil,
			}
		}

		err = newServer.Save()

		serverFolder := filepath.Join(serversFolder, serverId)

		err = copy.Copy(serverFolder, newFolders)
		if err != nil {
			fmt.Printf("Error migrating server %s files: %s\n", serverId, err)
			os.Exit(1)
		}
	}
}