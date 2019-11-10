/*
 Copyright 2018 Padduck, LLC

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

package shutdown

import (
	"github.com/braintree/manners"
	"github.com/pufferpanel/pufferpanel/v2/daemon/programs"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"os"
	"runtime/debug"
	"sync"
)

func CompleteShutdown() {
	logging.Info().Printf("Interrupt received, stopping servers")
	wg := Shutdown()
	wg.Wait()
	logging.Info().Printf("All servers stopped")
	os.Exit(0)
}

func Shutdown() *sync.WaitGroup {
	defer func() {
		if err := recover(); err != nil {
			logging.Error().Printf("%+v\n%s", err, debug.Stack())
		}
	}()
	wg := sync.WaitGroup{}
	programs.ShutdownService()
	manners.Close()
	prgs := programs.GetAll()
	wg.Add(len(prgs))
	for _, element := range prgs {
		go func(e *programs.Program) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					logging.Error().Printf("%+v\n%s", err, debug.Stack())
				}
			}()
			logging.Info().Printf("Stopping program %s", e.Id())
			running, err := e.IsRunning()
			if err != nil {
				logging.Error().Printf("Error stopping server %s: %s", e.Id(), err)
				return
			}
			if !running {
				return
			}
			err = e.Stop()
			if err != nil {
				logging.Error().Printf("Error stopping server %s: %s", e.Id(), err)
				return
			}
			err = e.GetEnvironment().WaitForMainProcess()
			if err != nil {
				logging.Error().Printf("Error stopping server %s: %s", e.Id(), err)
				return
			}
			logging.Info().Printf("Stopped program %s", e.Id())
		}(element)
	}
	return &wg
}
