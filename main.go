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
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferpanel/config"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/web"
	"os"
)

const Hash = "none"
const Version = "2.0.0-DEV"

func main() {
	r := gin.Default()
	web.RegisterRoutes(r)

	configPath, exists := os.LookupEnv("PUFFERPANEL_CONFIG")

	if !exists {
		configPath = "config.json"
	}

	cfg, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(configPath)
			if err != nil {
				logging.Error("Error creating config", err)
				return
			}
			data := config.Config{Database: config.Database{}}
			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			encoder.SetEscapeHTML(true)
			err = encoder.Encode(&data)
			if err != nil {
				logging.Error("Error creating config", err)
				return
			}
		} else {
			if err != nil {
				logging.Error("Error reading config", err)
				return
			}
			return
		}
	}

	err = config.Load(cfg)
	if err != nil {
		logging.Error("Error reading config", err)
		return
	}
	cfg.Close()

	logging.Init()

	err = database.Load()

	if err != nil {
		logging.Error("Error connecting to database", err)
	}

	defer database.Close()

	if err != nil {
		return
	}

	r.Run() // listen and serve on 0.0.0.0:8080
}