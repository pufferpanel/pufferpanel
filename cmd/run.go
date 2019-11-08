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
	"github.com/pufferpanel/apufferi/v4/logging"
	"github.com/pufferpanel/pufferd/v2/entry"
	"github.com/pufferpanel/pufferd/v2/sftp"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"github.com/pufferpanel/pufferpanel/v2/web"
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

func executeRun(cmd *cobra.Command, args []string) {
	err := internalRun(cmd, args)
	if err != nil {
		logging.Exception("An error has occurred while executing", err)
	}
}

func internalRun(cmd *cobra.Command, args []string) error {
	err := pufferpanel.LoadConfig()
	if err != nil {
		return err
	}

	err = logging.WithLogDirectory(viper.GetString("logs"), logging.DEBUG, nil)
	if err != nil {
		return err
	}

	logging.SetLevel(os.Stdout, logging.DEBUG)

	//load token, this also will store it to local node if there's one
	services.ValidateTokenLoaded()

	defer database.Close()

	services.LoadEmailService()

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithWriter(logging.AsWriter(logging.INFO)))

	web.RegisterRoutes(router)

	c := make(chan error)

	srv := &http.Server{
		Addr:    viper.GetString("web.host"),
		Handler: router,
	}

	httpsKey := viper.GetString("https.private")
	httpsCert := viper.GetString("https.public")

	if httpsKey != "" && httpsCert != "" {
		if _, err := os.Stat(httpsKey); err != nil {
			return err
		}

		if _, err := os.Stat(httpsCert); err != nil {
			return err
		}

		go func() {
			logging.Info("Listening for HTTPS requests on %s", srv.Addr)
			if err := srv.ListenAndServeTLS(httpsCert, httpsKey); err != nil && err != http.ErrServerClosed {
				c <- err
			}
		}()
	} else {
		go func() {
			logging.Info("Listening for HTTP requests on %s", srv.Addr)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				c <- err
			}
		}()
	}

	go func() {
		_, err := database.GetConnection()
		if err != nil {
			logging.Exception("Error connecting to database", err)
		}
	}()

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

	if viper.GetBool("localNode") {
		//local node!
		sftp.SetAuthorization(&services.DatabaseSFTPAuthorization{})
		c <- <-entry.Start()
	}

	return <-c
}
