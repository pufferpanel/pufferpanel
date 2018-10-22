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

package api

import (
	"github.com/gin-gonic/gin"
	builder "github.com/pufferpanel/apufferi/http"
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/pufferpanel/pufferpanel/models/view"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/shared"
	"net/http"
)

func registerUsers(g *gin.RouterGroup) {
	//g.Handle("GET", "", shared.NotImplemented)
	//g.Handle("OPTIONS", "", shared.CreateOptions("GET"))

	g.Handle("PUT", "/:username", CreateUser)
	g.Handle("GET", "/:username", GetUser)
	g.Handle("POST", "/:username", shared.NotImplemented)
	g.Handle("DELETE", "/:username", shared.NotImplemented)
	g.Handle("OPTIONS", "/:username", shared.CreateOptions("PUT", "GET", "POST", "DELETE"))
}

func CreateUser(c *gin.Context) {
	var us *services.UserService
	var err error
	response := builder.Respond(c)

	if us, err = services.GetUserService(); shared.HandleError(response, err) {
		return
	}

	var viewModel view.UserViewModel
	if err = c.BindJSON(&viewModel); shared.HandleError(response, err) {
		return
	}
	viewModel.Username = c.Param("username")

	user := &models.User{}
	viewModel.CopyToModel(user)

	if err = us.Create(user); shared.HandleError(response, err) {
		return
	}

	response.Data(view.FromUser(user)).Send()
}

func GetUser(c *gin.Context) {
	var us *services.UserService
	var err error
	response := builder.Respond(c)

	if us, err = services.GetUserService(); shared.HandleError(response, err) {
		return
	}

	username := c.Param("username")

	user, exists, err := us.Get(username)
	if shared.HandleError(response, err) {
		return
	} else if !exists {
		response.Fail().Status(http.StatusNotFound).Message("no user with username").Send()
		return
	}

	response.Data(view.FromUser(user)).Send()
}