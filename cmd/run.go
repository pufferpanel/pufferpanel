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
	"github.com/pufferpanel/pufferpanel/v2/daemon/entry"
	"github.com/pufferpanel/pufferpanel/v2/daemon/sftp"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/panel/database"
	"github.com/pufferpanel/pufferpanel/v2/panel/services"
	"github.com/pufferpanel/pufferpanel/v2/panel/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
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
		//load token, this also will store it to local node if there's one
		services.ValidateTokenLoaded()

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
	}

	if !noDaemon {
		go func() {
			c <- <-entry.Start()
		}()
	}

	return <-c
}
