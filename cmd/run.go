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

package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/common"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferpanel/config"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var runCmd = &cobra.Command{
	Use: "run",
	Short: "Runs the panel",
	Run: root,
}

func executeRun(cmd *cobra.Command, args []string) error {
	err := config.Load()
	if err != nil {
		return err
	}

	err = logging.WithLogDirectory("logs", logging.DEBUG, nil)
	if err != nil {
		return err
	}

	logging.SetLevel(os.Stdout, logging.DEBUG)

	err = database.Load()
	if err != nil {
		return err
	}

	defer database.Close()

	services.LoadEmailService()

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithWriter(logging.AsWriter(logging.INFO)))

	web.RegisterRoutes(router)

	c := make(chan error)

	srv := &http.Server{
		Addr: viper.GetString("web.host") + ":" + viper.GetString("web.port"),
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

	if runtime.GOOS == "linux" {
		go func() {
			file := viper.GetString("web.socket")

			if file == "" {
				return
			}

			err := os.Remove(file)
			if err != nil && !os.IsNotExist(err) {
				logging.Exception(fmt.Sprintf("Error deleting %s", file), err)
				return
			}

			listener, err := net.Listen("unix", file)
			defer common.Close(listener)
			if err != nil {
				logging.Exception(fmt.Sprintf("Error listening on %s", file), err)
				return
			}

			err = os.Chmod(file, 0777)
			if err != nil {
				logging.Exception(fmt.Sprintf("Error listening on %s", file), err)
				return
			}

			logging.Info("Listening for socket requests")
			c <- http.Serve(listener, router)
		}()
	}

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

	return <- c
}
