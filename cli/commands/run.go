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

package commands

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/cli"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/web"
)

type Run struct {
	cli.Command
	run bool
}

func (r *Run) Load() {
	flag.BoolVar(&r.run, "run", false, "Runs the service")
}

func (r *Run) ShouldRun() bool {
	return r.run
}

func (*Run) ShouldRunNext() bool {
	return false
}

func (r *Run) Run() error {
	err := logging.WithLogDirectory("logs", logging.DEBUG, nil)
	if err != nil {
		return err
	}

	err = database.Load()
	if err != nil {
		return err
	}

	defer database.Close()

	services.LoadEmailService()

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithWriter(logging.Writer))

	web.RegisterRoutes(router)

	return router.Run() // listen and serve on 0.0.0.0:8080
}
