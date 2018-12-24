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

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/web"
)

const Hash = "none"
const Version = "2.0.0-DEV"

func main() {
	logging.Init()
	logging.SetLevel(logging.DEBUG)

	err := database.Load()

	if err != nil {
		logging.Error("Error connecting to database", err)
	}

	defer database.Close()

	r := gin.Default()
	web.RegisterRoutes(r)

	err = r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		logging.Error("Error running web service", err)
	}
}