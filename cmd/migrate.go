package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/otiai10/copy"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/environments"
	"github.com/pufferpanel/pufferpanel/v2/legacy"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/programs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrates from v1 to v2",
	Run:   migrate,
}

var oldDaemonConfig string
var oldPanelConfig string

func init() {
	migrateCmd.Flags().StringVarP(&oldDaemonConfig, "daemon", "d", "/srv/puffferd/config.json", "Location of old pufferd config")
	migrateCmd.Flags().StringVarP(&oldPanelConfig, "panel", "p", "/etc/pufferpanel/config.json", "Location of old panel config")
}

func migrate(cmd *cobra.Command, args []string) {
	/*var confirm bool
	err := survey.AskOne(&survey.Confirm{
		Message: "Are you SURE you wish to migrate from v1 to v2? There is NO WARRANTY OR GUARANTEE this option will fully migrate your servers.",
		Default: false,
	}, &confirm)

	if err != nil {
		fmt.Printf("Error loading qusetion: %s\n", err)
		os.Exit(1)
		return
	}

	if !confirm {
		return
	}*/

	logging.DisableFileLogger()

	err := pufferpanel.LoadConfig("")
	if err != nil {
		fmt.Printf("Error loading new config: %s\n", err)
		os.Exit(1)
		return
	}

	if viper.GetBool("panel.enable") {
		migratePanel()
	}

	if viper.GetBool("daemon.enable") {
		migrateDaemon()
	}
}

func migratePanel() {
	oldConfig := &legacy.PanelConfig{}
	data, err := ioutil.ReadFile(oldPanelConfig)
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

	var newDb, oldDb, newDbConn *gorm.DB
	newDbConn, err = database.GetConnection()
	if err != nil {
		fmt.Printf("Error connection to new database: %s\n", err)
		os.Exit(1)
		return
	}

	newDb = newDbConn.Begin()
	defer func() {
		newDb.RollbackUnlessCommitted()
		newDbConn.Close()
	}()

	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", oldConfig.Mysql.Username, oldConfig.Mysql.Password, oldConfig.Mysql.Host, oldConfig.Mysql.Port, oldConfig.Mysql.Database)
	oldDb, err = gorm.Open("mysql", connString)
	if err != nil {
		fmt.Printf("Error connection to old database: %s\n", err)
		os.Exit(1)
		return
	}

	defer oldDb.Close()

	//migrate users
	var users []legacy.User
	err = oldDb.Find(&users).Error
	if err != nil {
		fmt.Printf("Error finding old users: %s\n", err)
		os.Exit(1)
		return
	}

	for _, v := range users {
		newUser := &models.User{
			ID:             v.ID,
			Username:       v.Username,
			Email:          v.Email,
			HashedPassword: v.Password,
		}
		err = newDb.Save(newUser).Error
		if err != nil {
			fmt.Printf("Error saving new user: %s\n", err)
			os.Exit(1)
			return
		}
	}

	//migrate nodes
	var nodes []legacy.Node
	err = oldDb.Find(&nodes).Error
	if err != nil {
		fmt.Printf("Error finding old nodes: %s\n", err)
		os.Exit(1)
		return
	}

	for _, v := range nodes {
		newNode := &models.Node{
			Name:        v.Name,
			PublicHost:  v.FQDN,
			PrivateHost: v.Ip,
			PublicPort:  v.Port,
			PrivatePort: v.Port,
			SFTPPort:    v.Sftp,
			Secret:      v.DaemonSecret,
		}
		err = newDb.Save(newNode).Error
		if err != nil {
			fmt.Printf("Error saving new node: %s\n", err)
			os.Exit(1)
			return
		}
	}

	//migrate servers
	var servers []legacy.Server
	err = oldDb.Find(&servers).Error
	if err != nil {
		fmt.Printf("Error finding old nodes: %s\n", err)
		os.Exit(1)
		return
	}

	for _, v := range servers {
		newServer := &models.Server{
			Name:       v.Name,
			Identifier: v.Hash.String(),
			NodeID:     v.Node,
			Node:       models.Node{},
			IP:         "?",
			Port:       0,
			Type:       "migrated",
		}
		err = newDb.Save(newServer).Error
		if err != nil {
			fmt.Printf("Error saving new server: %s\n", err)
			os.Exit(1)
			return
		}

		//servers in the old system used an owner column, new one does not, we use oauth
		//so.... we need to grant them
		perms := &models.Permissions{
			ServerIdentifier: &newServer.Identifier,
			UserId: &v.OwnerId,
		}
		perms.SetDefaults()

		err = newDb.Save(perms).Error
		if err != nil {
			fmt.Printf("Error saving new permission: %s\n", err)
			os.Exit(1)
			return
		}
	}

	//migrate permissions

	//migrate clients

	newDb.Commit()
}

func migrateDaemon() {
	oldConfig := &legacy.Config{}
	data, err := ioutil.ReadFile(oldDaemonConfig)
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

	programs.ServerFolder = viper.GetString("daemon.data.servers")

	environments.LoadModules()

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
			Server: pufferpanel.Server{
				Variables:      map[string]pufferpanel.Variable{},
				Display:        legacyData.ProgramData.Display,
				Environment:    legacyData.ProgramData.EnvironmentData,
				Installation:   convertCommands(legacyData.ProgramData.InstallData.Operations),
				Uninstallation: convertCommands(legacyData.ProgramData.UninstallData.Operations),
				Type:           pufferpanel.Type{Type: legacyData.ProgramData.Type},
				Identifier:     serverId,
				Execution: pufferpanel.Execution{
					Command:                 strings.TrimSpace(legacyData.ProgramData.RunData.Program + " " + strings.Join(legacyData.ProgramData.RunData.Arguments, " ")),
					StopCommand:             legacyData.ProgramData.RunData.Stop,
					Disabled:                !legacyData.ProgramData.RunData.Enabled,
					AutoStart:               legacyData.ProgramData.RunData.AutoStart,
					AutoRestartFromCrash:    legacyData.ProgramData.RunData.AutoRestartFromCrash,
					AutoRestartFromGraceful: legacyData.ProgramData.RunData.AutoRestartFromGraceful,
					PreExecution:            convertCommands(legacyData.ProgramData.RunData.Pre),
					PostExecution:           nil,
					StopCode:                legacyData.ProgramData.RunData.StopCode,
					EnvironmentVariables:    legacyData.ProgramData.RunData.EnvironmentVariables,
				},
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

		err = copy.Copy(serverFolder, filepath.Join(programs.ServerFolder, serverId))
		if err != nil {
			fmt.Printf("Error migrating server %s files: %s\n", serverId, err)
			os.Exit(1)
		}
	}
}

func convertCommands(source []map[string]interface{}) []interface{} {
	if source == nil {
		return nil
	}
	result := make([]interface{}, len(source))
	for k, v := range source {
		result[k] = v
	}
	return result
}
