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
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/daemon"
	"github.com/pufferpanel/pufferpanel/v2/daemon/entry"
	"github.com/pufferpanel/pufferpanel/v2/daemon/sftp"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/panel/database"
	"github.com/pufferpanel/pufferpanel/v2/panel/models"
	"github.com/pufferpanel/pufferpanel/v2/panel/services"
	"github.com/pufferpanel/pufferpanel/v2/panel/web"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
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

var noWeb bool
var noDaemon bool

func init() {
	runCmd.Flags().BoolVar(&noWeb, "noweb", false, "Do not run web interface")
	runCmd.Flags().BoolVar(&noDaemon, "nodaemon", false, "Do not run daemon")
}

func executeRun(cmd *cobra.Command, args []string) {
	err := internalRun(cmd, args)
	if err != nil {
		logging.Error().Printf("An error occurred: %s", err.Error())
	}
}

func internalRun(cmd *cobra.Command, args []string) error {
	err := pufferpanel.LoadConfig("")
	if err != nil {
		return err
	}

	defer logging.Close()

	c := make(chan error)

	go func() {
		quit := make(chan os.Signal)
		// kill (no param) default send syscall.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall.SIGKILL but can"t be catch, so don't need add it
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down web server ...")
		c <- nil
	}()

	if !noWeb {
		services.ValidateTokenLoaded()
		daemon.SetPublicKey(services.GetPublicKey())

		defer database.Close()

		services.LoadEmailService()

		router := gin.New()
		router.Use(gin.Recovery())
		router.Use(gin.LoggerWithWriter(logging.Debug().Writer()))

		web.RegisterRoutes(router)

		srv := &http.Server{
			Addr:    viper.GetString("panel.web.host"),
			Handler: router,
		}

		go func() {
			logging.Info().Printf("Listening for HTTP requests on %s", srv.Addr)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				c <- err
			}
		}()

		go func() {
			_, err := database.GetConnection()
			if err != nil {
				logging.Error().Printf("Error connecting to database: %s", err.Error())
			}
		}()

		//if we have the web, then let's use our sftp auth instead
		sftp.SetAuthorization(&services.DatabaseSFTPAuthorization{})

		//validate local daemon is configured in this panel
		if !noDaemon {
			go func() {
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
						PublicPort:  5656,
						PrivatePort: 5656,
						SFTPPort:    5657,
						Secret:      strings.Replace(uuid.NewV4().String(), "-", "", -1),
					}
					nodeHost := viper.GetString("daemon.web.host")
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
			}()
		}
	}

	if !noDaemon {
		go func() {
			c <- <-entry.Start()
		}()
	}

	return <-c
}
