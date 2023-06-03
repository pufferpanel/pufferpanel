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
	"github.com/pufferpanel/pufferpanel/v2/programs"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"github.com/pufferpanel/pufferpanel/v2/sftp"
	"github.com/pufferpanel/pufferpanel/v2/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the panel",
	Run:   executeRun,
}

var webService *manners.GracefulServer

func executeRun(cmd *cobra.Command, args []string) {
	term := make(chan bool, 10)

	internalRun(term)
	//wait for the termination signal, so we can shut down
	<-term

	//shut down everything
	//all of these can be closed regardless of what type of install this is, as they all check if they are even being
	//used
	logging.Debug.Printf("stopping http server")
	if webService != nil {
		webService.Close()
	}

	logging.Debug.Printf("stopping sftp server")
	sftp.Stop()

	logging.Debug.Printf("stopping servers")
	programs.ShutdownService()
	for _, p := range programs.GetAll() {
		_ = p.Stop()
		p.RunningEnvironment.WaitForMainProcessFor(time.Minute) //wait 60 seconds
	}

	logging.Debug.Printf("stopping database connections")
	database.Close()
}

func internalRun(terminate chan bool) {
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

	if config.PanelEnabled.Value() {
		panel()

		if config.SessionKey.Value() == "" {
			k := securecookie.GenerateRandomKey(32)
			if err := config.SessionKey.Set(hex.EncodeToString(k), true); err != nil {
				logging.Error.Printf("error saving session key: %s", err.Error())
				terminate <- true
				return
			}
		}

		result, err := hex.DecodeString(config.SessionKey.Value())
		if err != nil {
			logging.Error.Printf("error decoding session key: %s", err.Error())
			terminate <- true
			return
		}
		sessionStore := cookie.NewStore(result)
		router.Use(sessions.Sessions("session", sessionStore))
	}

	if config.DaemonEnabled.Value() {
		err := daemon()
		if err != nil {
			logging.Error.Printf("error starting daemon server: %s", err.Error())
			terminate <- true
			return
		}
	}

	web.RegisterRoutes(router)

	go func() {
		l, err := net.Listen("tcp", config.WebHost.Value())
		if err != nil {
			logging.Error.Printf("error starting http server: %s", err.Error())
			terminate <- true
			return
		}

		logging.Info.Printf("Listening for HTTP requests on %s", l.Addr().String())
		webService = manners.NewWithServer(&http.Server{Handler: router})
		err = webService.Serve(l)
		if err != nil && err != http.ErrServerClosed {
			logging.Error.Printf("error listening for http requests: %s", err.Error())
			terminate <- true
		}
	}()

	return
}

func panel() {
	services.LoadEmailService()

	//if we have the web, then let's use our sftp auth instead
	sftp.SetAuthorization(&services.DatabaseSFTPAuthorization{})
}

func daemon() error {
	sftp.Run()

	environments.LoadModules()
	programs.Initialize()

	var err error

	if _, err = os.Stat(config.ServersFolder.Value()); os.IsNotExist(err) {
		logging.Info.Printf("No server directory found, creating")
		err = os.MkdirAll(config.ServersFolder.Value(), 0755)
		if err != nil && !os.IsExist(err) {
			return err
		}
	}

	//check if viper directly has a value for binaries, so we can migrate
	if runtime.GOOS == "linux" && !viper.InConfig(config.BinariesFolder.Key()) {
		_ = config.BinariesFolder.Set(filepath.Join(filepath.Dir(config.ServersFolder.Value()), "binaries"), true)
	}

	err = os.MkdirAll(config.BinariesFolder.Value(), 0755)
	if err != nil {
		logging.Error.Printf("Error creating binaries folder: %s", err.Error())
	}

	//update path to include our binary folder
	newPath := os.Getenv("PATH")
	fullPath, _ := filepath.Abs(config.BinariesFolder.Value())
	if !strings.Contains(newPath, fullPath) {
		_ = os.Setenv("PATH", newPath+":"+fullPath)
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
	return nil
}
