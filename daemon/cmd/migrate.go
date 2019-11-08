package main

import (
	"encoding/json"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v2/daemon"
	"github.com/pufferpanel/pufferpanel/v2/daemon/programs"
	"github.com/pufferpanel/pufferpanel/v2/shared"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "This will migrate v1 servers to v2",
	Run:   migrate,
}

func migrate(cmd *cobra.Command, args [] string) {
	_ = daemon.LoadConfig()

	programFiles, err := ioutil.ReadDir(programs.ServerFolder)
	if err != nil {
		fmt.Printf("Error reading from server data folder: %s\n", err)
		return
	}

	backupDir := path.Join(programs.ServerFolder, "backup")

	err = os.MkdirAll(backupDir, 0755)
	if err != nil {
		fmt.Printf("Error creating backup directory %s: %s\n", backupDir, err)
		return
	}

	for _, element := range programFiles {
		if element.IsDir() || !strings.HasSuffix(element.Name(), ".json") {
			continue
		}

		var fullPath = path.Join(programs.ServerFolder, element.Name())

		err = backupFile(fullPath, backupDir)
		if err != nil {
			fmt.Printf("Error backing up file: %s\n", err)
			continue
		}

		err = migrateFile(fullPath)
		if err != nil {
			fmt.Printf("Error migrating file: %s\n", err)
			continue
		}
	}
}

func migrateFile(name string) (err error) {
	fmt.Printf("Attempting to migrate %s\n", name)

	program, err := loadFile(name)
	if err != nil {
		return
	}

	//copy base items over that map directly
	replacement := shared.Server{
		Variables: make(map[string]shared.Variable),
		Display:   program.ProgramData.Display,
		Environment: shared.Type{
			Type: program.ProgramData.EnvironmentData["type"].(string),
		},
		Installation:   make([]interface{}, 0),
		Uninstallation: make([]interface{}, 0),
		Type:           program.ProgramData.Type,
		Identifier:     program.ProgramData.Identifier,
		Execution: shared.Execution{
			Arguments:   program.ProgramData.RunData.Arguments,
			ProgramName: program.ProgramData.RunData.Program,
			StopCommand: program.ProgramData.RunData.Stop,
			//Disabled:                !program.ProgramData.RunData.Enabled,
			Disabled:                false,
			AutoStart:               program.ProgramData.RunData.AutoStart,
			AutoRestartFromCrash:    program.ProgramData.RunData.AutoRestartFromCrash,
			AutoRestartFromGraceful: program.ProgramData.RunData.AutoRestartFromGraceful,
			PreExecution:            make([]interface{}, 0),
			PostExecution:           make([]interface{}, 0),
			StopCode:                program.ProgramData.RunData.StopCode,
			EnvironmentVariables:    program.ProgramData.RunData.EnvironmentVariables,
		},
	}

	//copy data
	for k, v := range program.ProgramData.Data {
		replacement.Variables[k] = shared.Variable{
			Description:  v.Description,
			Display:      v.Display,
			Internal:     v.Internal,
			Required:     v.Required,
			Value:        v.Value,
			UserEditable: v.UserEditable,
			Type:         "text",
			Options:      nil,
		}
	}

	//copy installation and uninstall sections by removing type from map
	for _, v := range program.ProgramData.InstallData.Operations {
		c := make(map[string]interface{})
		for i, o := range v {
			if i == "type" {
				continue
			}
			c[i] = o
		}

		replacement.Installation = append(replacement.Installation, shared.MetadataType{
			Type:     v["type"].(string),
			Metadata: c,
		})
	}

	for _, v := range program.ProgramData.UninstallData.Operations {
		c := make(map[string]interface{})
		for i, o := range v {
			if i == "type" {
				continue
			}
			c[i] = o
		}

		replacement.Uninstallation = append(replacement.Uninstallation, shared.MetadataType{
			Type:     v["type"].(string),
			Metadata: c,
		})
	}

	for _, v := range program.ProgramData.RunData.Pre {
		c := make(map[string]interface{})
		for i, o := range v {
			if i == "type" {
				continue
			}
			c[i] = o
		}

		replacement.Execution.PreExecution = append(replacement.Execution.PreExecution, shared.MetadataType{
			Type:     v["type"].(string),
			Metadata: c,
		})
	}

	for _, v := range program.ProgramData.RunData.Post {
		c := make(map[string]interface{})
		for i, o := range v {
			if i == "type" {
				continue
			}
			c[i] = o
		}

		replacement.Execution.PostExecution = append(replacement.Execution.PostExecution, shared.MetadataType{
			Type:     v["type"].(string),
			Metadata: c,
		})
	}

	target, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	encoder := json.NewEncoder(target)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(replacement)
	return
}

func backupFile(name string, targetDir string) error {
	fmt.Printf("Attempting to backup %s to %s\n", name, targetDir)
	source, err := os.Open(name)
	if err != nil {
		return err
	}
	defer shared.Close(source)

	filename := path.Base(name)

	target, err := os.OpenFile(path.Join(targetDir, filename), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer shared.Close(target)

	_, err = io.Copy(target, source)
	return err
}

func loadFile(name string) (program ServerJson, err error) {
	file, err := os.Open(name)
	if err != nil {
		return
	}
	defer shared.Close(file)

	err = json.NewDecoder(file).Decode(&program)
	return
}

//copy of the old v1 objects
type ServerJson struct {
	ProgramData ProgramData `json:"pufferd"`
}

type ProgramData struct {
	Data            map[string]DataObject  `json:"data,omitempty"`
	Display         string                 `json:"display,omitempty"`
	EnvironmentData map[string]interface{} `json:"environment,omitempty"`
	InstallData     InstallSection         `json:"install,omitempty"`
	UninstallData   InstallSection         `json:"uninstall,omitempty"`
	Type            string                 `json:"type,omitempty"`
	Identifier      string                 `json:"id,omitempty"`
	RunData         RunObject              `json:"run,omitempty"`
	Template        string                 `json:"template,omitempty"`
}

type DataObject struct {
	Description  string      `json:"desc,omitempty"`
	Display      string      `json:"display,omitempty"`
	Internal     bool        `json:"internal,omitempty"`
	Required     bool        `json:"required,omitempty"`
	Value        interface{} `json:"value,omitempty"`
	UserEditable bool        `json:"userEdit,omitempty"`
}

type RunObject struct {
	Arguments               []string                 `json:"arguments,omitempty"`
	Program                 string                   `json:"program,omitempty"`
	Stop                    string                   `json:"stop,omitempty"`
	Enabled                 bool                     `json:"enabled,omitempty"`
	AutoStart               bool                     `json:"autostart,omitempty"`
	AutoRestartFromCrash    bool                     `json:"autorecover,omitempty"`
	AutoRestartFromGraceful bool                     `json:"autorestart,omitempty"`
	Pre                     []map[string]interface{} `json:"pre,omitempty"`
	Post                    []map[string]interface{} `json:"post,omitempty"`
	StopCode                int                      `json:"stopCode,omitempty"`
	EnvironmentVariables    map[string]string        `json:"environmentVars,omitempty"`
}

type InstallSection struct {
	Operations []map[string]interface{} `json:"commands,,omitempty"`
}
