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
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/web"
	"os"
)

const Hash = "none"
const Version = "2.0.0-DEV"

func main() {
	logging.Init()
	logging.SetLevel(logging.DEBUG)

	err := database.Load()

	if err != nil {
		logging.Error("Error connecting to database", err)
		return
	}

	defer database.Close()

	args := os.Args[1:]

	run := false

	if len(args) > 0 {
		counter := 0
		for counter < len(args) {
			arg := args[counter]
			switch arg {
			case "--addUser":
				{
					if counter+3 >= len(args) {
						logging.Errorf("not enough arguments to create user")
						return
					}
					username := args[counter+1]
					email := args[counter+2]
					password := args[counter+3]
					counter += 3

					us, err := services.GetUserService()
					if err != nil {
						logging.Error("could not load user service", err)
						return
					}
					user := &models.User{Email: email, Username: username}
					err = user.SetPassword(password)

					if err != nil {
						logging.Error("could not create user", err)
						return
					}

					err = us.Create(user)
					if err != nil {
						logging.Error("could not create user", err)
						return
					}
				}
			case "--run", "-r":
				{
					run = true
				}
			default:
				{
					logging.Warnf("unknown argument: %s", arg)
				}
			}
			counter++
		}
	}

	if !run {
		return
	}

	services.LoadEmailService()

	r := gin.Default()
	web.RegisterRoutes(r)

	err = r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		logging.Error("Error running web service", err)
	}
}

func createUser(username, email, password string) {

}
