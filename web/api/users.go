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
	"strconv"
)

const MAX_PAGE_SIZE = 100
const DEFAULT_PAGE_SIZE = 20

func registerUsers(g *gin.RouterGroup) {
	g.Handle("GET", "", SearchUsers)
	g.Handle("OPTIONS", "", shared.CreateOptions("GET"))

	g.Handle("PUT", "/:username", CreateUser)
	g.Handle("GET", "/:username", GetUser)
	g.Handle("POST", "/:username", UpdateUser)
	g.Handle("DELETE", "/:username", DeleteUser)
	g.Handle("OPTIONS", "/:username", shared.CreateOptions("PUT", "GET", "POST", "DELETE"))
}

func SearchUsers (c *gin.Context) {
	var us *services.UserService
	var err error
	response := builder.Respond(c)

	usernameFilter := c.DefaultQuery("username", "*")
	emailFilter := c.DefaultQuery("email", "*")
	pageSizeQuery := c.DefaultQuery("limit", strconv.Itoa(DEFAULT_PAGE_SIZE))
	pageQuery := c.DefaultQuery("page", strconv.Itoa(1))

	pageSize, err := strconv.Atoi(pageSizeQuery)
	if err != nil || pageSize <= 0 {
		response.Fail().Status(http.StatusBadRequest).Message("page size must be a positive number").Send()
		return
	}

	if pageSize > MAX_PAGE_SIZE {
		pageSize = MAX_PAGE_SIZE
	}

	page, err := strconv.Atoi(pageQuery)
	if err != nil || page <= 0 {
		response.Fail().Status(http.StatusBadRequest).Message("page must be a positive number").Send()
		return
	}

	if us, err = services.GetUserService(); shared.HandleError(response, err) {
		return
	}

	var results *models.Users
	if results, err = us.Search(usernameFilter, emailFilter, uint(pageSize), uint(page)); shared.HandleError(response, err) {
		return
	}

	response.Data(view.FromUsers(results)).Send()
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

func UpdateUser(c *gin.Context) {
	var us *services.UserService
	var err error
	response := builder.Respond(c)

	if us, err = services.GetUserService(); shared.HandleError(response, err) {
		return
	}

	username := c.Param("username")

	var viewModel view.UserViewModel
	if err = c.BindJSON(&viewModel); shared.HandleError(response, err) {
		return
	}

	user, exists, err := us.Get(username)
	if shared.HandleError(response, err) {
		return
	} else if !exists {
		response.Fail().Status(http.StatusNotFound).Message("no user with username").Send()
		return
	}

	viewModel.CopyToModel(user)
	us.Update(user)

	response.Data(view.FromUser(user)).Send()
}

func DeleteUser(c *gin.Context) {
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

	if err = us.Delete(user.Username); shared.HandleError(response, err) {
		return
	}

	response.Data(view.FromUser(user)).Send()
}