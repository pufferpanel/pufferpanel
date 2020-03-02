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

package daemon

import (
	"fmt"
	"github.com/braintree/manners"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v2/daemon/web"
	"github.com/pufferpanel/pufferpanel/v2/environments"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/programs"
	"github.com/pufferpanel/pufferpanel/v2/sftp"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

var runService = true

var Engine *gin.Engine

func Start() chan error {
	errChan := make(chan error)
	go entry(errChan)
	return errChan
}

func entry(errChan chan error) {
	environments.LoadModules()
	programs.Initialize()

	var err error

	if _, err = os.Stat(programs.ServerFolder); os.IsNotExist(err) {
		logging.Error().Printf("No server directory found, creating")
		err = os.MkdirAll(programs.ServerFolder, 0755)
		if err != nil && !os.IsExist(err) {
			errChan <- err
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

	createHook()

	for runService && err == nil {
		err = runServices()
	}

	Shutdown()

	errChan <- err
	return
}

func runServices() error {
	defer recoverPanic()

	Engine = web.ConfigureWeb()

	sftp.Run()

	web := viper.GetString("daemon.web.host")

	logging.Debug().Printf("Starting web access on %s", web)
	err := manners.ListenAndServe(web, Engine)

	return err
}

func createHook() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGPIPE)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logging.Error().Printf("%+v\n%s", err, debug.Stack())
			}
		}()

		var sig os.Signal

		for sig != syscall.SIGTERM {
			sig = <-c
			switch sig {
			case syscall.SIGPIPE:
				//ignore SIGPIPEs for now, we're somehow getting them and it's causing issues
			}
		}

		runService = false
		CompleteShutdown()
	}()
}

func recoverPanic() {
	if rec := recover(); rec != nil {
		err := rec.(error)
		fmt.Printf("CRITICAL: %s", err.Error())
		logging.Error().Printf("Unhandled error: %s", err.Error())
	}
}
