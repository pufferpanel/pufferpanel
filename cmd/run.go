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
	"encoding/hex"
	"github.com/braintree/manners"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/config"
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
	c := make(chan error)
	term := make(chan bool)

	internalRun(c, term)
	err := <-c
	if err != nil {
		logging.Error.Printf("An error occurred: %s", err.Error())
	}
}

func internalRun(c chan error, terminate chan bool) {
	if err := config.LoadConfigFile(""); err != nil {
		c <- err
		return
	}

	logging.Initialize(true)
	signal.Ignore(syscall.SIGPIPE, syscall.SIGHUP)

	go func() {
		quit := make(chan os.Signal)
		// kill (no param) default send syscall.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall.SIGKILL but can"t be catch, so don't need add it
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logging.Info.Println("Shutting down...")
		terminate <- true
	}()

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithWriter(logging.Info.Writer()))
	gin.DefaultWriter = logging.Info.Writer()
	gin.DefaultErrorWriter = logging.Error.Writer()
	pufferpanel.Engine = router

	if config.GetBool("panel.enable") {
		panel(c)

		if config.GetString("panel.sessionKey") == "" {
			if err := config.Set("panel.sessionKey", securecookie.GenerateRandomKey(32)); err != nil {
				c <- err
				return
			}
		}

		result, err := hex.DecodeString(config.GetString("panel.sessionKey"))
		if err != nil {
			c <- err
			return
		}
		sessionStore := cookie.NewStore(result)
		router.Use(sessions.Sessions("session", sessionStore))
	}

	if config.GetBool("daemon.enable") {
		daemon(c)
	}

	web.RegisterRoutes(router)

	go func() {
		l, err := net.Listen("tcp", config.GetString("web.host"))
		if err != nil {
			c <- err
			return
		}

		logging.Info.Printf("Listening for HTTP requests on %s", l.Addr().String())
		err = manners.Serve(l, router)
		if err != nil && err != http.ErrServerClosed {
			c <- err
			terminate <- true
		}
	}()

	go func() {
		//wait for the termination signal, so we can shut down
		<-terminate

		//shut down everything
		//all of these can be closed regardless of what type of install this is, as they all check if they are even being
		//used

		manners.Close()
		sftp.Stop()
		programs.ShutdownService()
		database.Close()

		//return out, the upper layers know how to handle this
		c <- nil
	}()

	return
}

func panel(ch chan error) {
	services.ValidateTokenLoaded()
	services.LoadEmailService()

	//if we have the web, then let's use our sftp auth instead
	sftp.SetAuthorization(&services.DatabaseSFTPAuthorization{})

	err := config.LoadConfigDatabase(database.GetConnector())
	if err != nil {
		logging.Error.Printf("Error loading config from database: %s", err.Error())
	}

	//validate local daemon is configured in this panel
	if config.GetBool("daemon.enable") {
		db, err := database.GetConnection()
		if err != nil {
			return
		}
		ns := &services.Node{DB: db}
		nodes, err := ns.GetAll()
		if err != nil {
			logging.Error.Printf("Failed to get nodes: %s", err.Error())
			return
		}

		if len(*nodes) == 0 {
			logging.Info.Printf("Adding local node")
			create := &models.Node{
				Name:        "LocalNode",
				PublicHost:  "127.0.0.1",
				PrivateHost: "127.0.0.1",
				PublicPort:  8080,
				PrivatePort: 8080,
				SFTPPort:    5657,
				Secret:      strings.Replace(uuid.NewV4().String(), "-", "", -1),
			}
			nodeHost := config.GetString("web.host")
			sftpHost := config.GetString("daemon.sftp.host")
			hostParts := strings.SplitN(nodeHost, ":", 2)
			sftpParts := strings.SplitN(sftpHost, ":", 2)

			if len(hostParts) == 2 {
				port, err := strconv.Atoi(hostParts[1])
				if err == nil {
					create.PublicPort = uint16(port)
					create.PrivatePort = uint16(port)
				}
			}
			if len(sftpParts) == 2 {
				port, err := strconv.Atoi(sftpParts[1])
				if err == nil {
					create.SFTPPort = uint16(port)
				}
			}

			err = ns.Create(create)
			if err != nil {
				logging.Error.Printf("Failed to add local node: %s", err.Error())
			}
		}
	}
}

func daemon(ch chan error) {
	sftp.Run()

	pufferpanel.InitEnvironment()
	environments.LoadModules()
	programs.Initialize()

	var err error

	if _, err = os.Stat(pufferpanel.ServerFolder); os.IsNotExist(err) {
		logging.Error.Printf("No server directory found, creating")
		err = os.MkdirAll(pufferpanel.ServerFolder, 0755)
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
				logging.Info.Printf("Queued server %s", element.Id())
				element.GetEnvironment().DisplayToConsole(true, "Server has been queued to start\n")
				programs.StartViaService(element)
			}
		}
	}
}
