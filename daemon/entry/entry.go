package entry

import (
	"fmt"
	"github.com/braintree/manners"
	"github.com/pufferpanel/pufferpanel/v2/daemon/environments"
	"github.com/pufferpanel/pufferpanel/v2/daemon/programs"
	"github.com/pufferpanel/pufferpanel/v2/daemon/routing"
	"github.com/pufferpanel/pufferpanel/v2/daemon/sftp"
	"github.com/pufferpanel/pufferpanel/v2/daemon/shutdown"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

var runService = true

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

	shutdown.Shutdown()

	errChan <- err
	return
}

func runServices() error {
	defer recoverPanic()

	router := routing.ConfigureWeb()

	sftp.Run()

	web := viper.GetString("listen.web")

	logging.Debug().Printf("Starting web access on %s", web)
	err := manners.ListenAndServe(web, router)

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
		shutdown.CompleteShutdown()
	}()
}

func recoverPanic() {
	if rec := recover(); rec != nil {
		err := rec.(error)
		fmt.Printf("CRITICAL: %s", err.Error())
		logging.Error().Printf("Unhandled error: %s", err.Error())
	}
}
