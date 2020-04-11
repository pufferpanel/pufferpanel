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

package main

import (
	"github.com/braintree/manners"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/environments"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/programs"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"github.com/pufferpanel/pufferpanel/v2/sftp"
	"github.com/pufferpanel/pufferpanel/v2/web"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the panel",
	Run:   executeRun,
}

func executeRun(cmd *cobra.Command, args []string) {
	err := internalRun(cmd, args)
	if err != nil {
		logging.Error().Printf("An error occurred: %s", err.Error())
	}
}

func internalRun(cmd *cobra.Command, args []string) error {
	if err := pufferpanel.LoadConfig(""); err != nil {
		return err
	}

	logging.Initialize()

	defer logging.Close()
	defer database.Close()

	c := make(chan error)

	signal.Ignore(syscall.SIGPIPE, syscall.SIGHUP)

	go func() {
		quit := make(chan os.Signal)
		// kill (no param) default send syscall.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall.SIGKILL but can"t be catch, so don't need add it
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down...")
		c <- nil
	}()

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithWriter(logging.Debug().Writer()))
	pufferpanel.Engine = router

	web.RegisterRoutes(router)

	if viper.GetBool("panel.enable") {
		panel(c)
	}

	if viper.GetBool("daemon.enable") {
		daemon(c)
	}

	go func() {
		l, err := net.Listen("tcp4", viper.GetString("web.host"))
		if err != nil {
			c <- err
			return
		}

		logging.Info().Printf("Listening for HTTP requests on %s", l.Addr().String())
		for err = manners.Serve(l, router); err != nil && err != http.ErrServerClosed; err = manners.Serve(l, router) {
			c <- err
		}
	}()

	return <-c
}

func panel(ch chan error) {
	services.ValidateTokenLoaded()
	services.LoadEmailService()

	go func() {
		_, err := database.GetConnection()
		if err != nil {
			logging.Error().Printf("Error connecting to database: %s", err.Error())
		}
	}()

	//if we have the web, then let's use our sftp auth instead
	sftp.SetAuthorization(&services.DatabaseSFTPAuthorization{})

	//validate local daemon is configured in this panel
	if viper.GetBool("daemon.enable") {
		db, err := database.GetConnection()
		if err != nil {
			return
		}
		ns := &services.Node{DB: db}
		nodes, err := ns.GetAll()
		if err != nil {
			logging.Error().Printf("Failed to get nodes: %s", err.Error())
			return
		}
		exists := false
		for _, n := range *nodes {
			if n.IsLocal() {
				exists = true
			}
		}
		if !exists {
			logging.Info().Printf("Adding local node")
			create := &models.Node{
				Name:        "LocalNode",
				PublicHost:  "127.0.0.1",
				PrivateHost: "127.0.0.1",
				PublicPort:  8080,
				PrivatePort: 8080,
				SFTPPort:    5657,
				Secret:      strings.Replace(uuid.NewV4().String(), "-", "", -1),
			}
			nodeHost := viper.GetString("web.host")
			sftpHost := viper.GetString("daemon.sftp.host")
			hostParts := strings.SplitN(nodeHost, ":", 2)
			sftpParts := strings.SplitN(sftpHost, ":", 2)

			if len(hostParts) == 2 {
				port, err := strconv.Atoi(hostParts[1])
				if err == nil {
					create.PublicPort = uint(port)
					create.PrivatePort = uint(port)
				}
			}
			if len(sftpParts) == 2 {
				port, err := strconv.Atoi(sftpParts[1])
				if err == nil {
					create.SFTPPort = uint(port)
				}
			}

			//override ENV because we need our id here instead since it's new
			viper.Set("PUFFER_DAEMON_AUTH_CLIENTID", create.Secret)
			err = ns.Create(create)
			if err != nil {
				logging.Error().Printf("Failed to add local node: %s", err.Error())
			}
		}
	}
}

func daemon(ch chan error) {
	sftp.Run()

	environments.LoadModules()
	programs.Initialize()

	var err error

	if _, err = os.Stat(programs.ServerFolder); os.IsNotExist(err) {
		logging.Error().Printf("No server directory found, creating")
		err = os.MkdirAll(programs.ServerFolder, 0755)
		if err != nil && !os.IsExist(err) {
			ch <- err
			return
		}
	}

	programs.LoadFromFolder()

	programs.InitService()

	for _, element := range programs.GetAll() {
		if element.IsEnabled() {
			element.GetEnvironment().DisplayToConsole(true, "Daemon has been started\n")
			if element.IsAutoStart() {
				logging.Info().Printf("Queued server %s", element.Id())
				element.GetEnvironment().DisplayToConsole(true, "Server has been queued to start\n")
				programs.StartViaService(element)
			}
		}
	}
}
