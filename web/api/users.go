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
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/apufferi/v3/response"
	"github.com/pufferpanel/apufferi/v3/scope"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"github.com/pufferpanel/pufferpanel/v2/web/handlers"
	"net/http"
)

func registerUsers(g *gin.RouterGroup) {
	//if you can log in, you can see and edit yourself
	g.Handle("GET", "", handlers.OAuth2(scope.Login, false), getSelf)
	g.Handle("PUT", "", handlers.OAuth2(scope.Login, false), updateSelf)
	g.Handle("POST", "", handlers.OAuth2(scope.UsersView, false), searchUsers)
	g.Handle("OPTIONS", "", response.CreateOptions("GET", "PUT", "POST"))

	g.Handle("PUT", "/:username", handlers.OAuth2(scope.UsersEdit, false), createUser)
	g.Handle("GET", "/:username", handlers.OAuth2(scope.UsersView, false), getUser)
	g.Handle("POST", "/:username", handlers.OAuth2(scope.UsersEdit, false), updateUser)
	g.Handle("DELETE", "/:username", handlers.OAuth2(scope.UsersEdit, false), deleteUser)
	g.Handle("OPTIONS", "/:username", response.CreateOptions("PUT", "GET", "POST", "DELETE"))
}

func searchUsers(c *gin.Context) {
	var err error
	res := response.From(c)

	search := newUserSearch()
	err = c.BindJSON(search)
	if response.HandleError(res, err) {
		return
	}
	if search.PageLimit <= 0 {
		res.Fail().Status(http.StatusBadRequest).Message("page size must be a positive number")
		return
	}

	if search.PageLimit > MaxPageSize {
		search.PageLimit = MaxPageSize
	}

	if search.Page <= 0 {
		res.Fail().Status(http.StatusBadRequest).Message("page must be a positive number")
		return
	}

	db, err := database.GetConnection()
	if response.HandleError(res, err) {
		return
	}

	us := &services.User{DB: db}

	var results *models.Users
	var total uint
	if results, total, err = us.Search(search.Username, search.Email, uint(search.PageLimit), uint(search.Page)); response.HandleError(res, err) {
		return
	}

	res.PageInfo(uint(search.Page), uint(search.PageLimit), MaxPageSize, total).Data(models.FromUsers(results))
}

func createUser(c *gin.Context) {
	var err error
	res := response.Respond(c)

	db, err := database.GetConnection()
	if response.HandleError(res, err) {
		return
	}

	us := &services.User{DB: db}

	var viewModel models.UserView
	if err = c.BindJSON(&viewModel); response.HandleError(res, err) {
		return
	}
	viewModel.Username = c.Param("username")

	if err = viewModel.Valid(false); response.HandleError(res, err) {
		return
	}

	if viewModel.Password == "" {
		response.HandleError(res, pufferpanel.ErrFieldRequired("password"))
		return
	}

	user := &models.User{}
	viewModel.CopyToModel(user)

	if err = us.Create(user); response.HandleError(res, err) {
		return
	}

	res.Data(models.FromUser(user))
}

func getUser(c *gin.Context) {
	res := response.From(c)

	db, err := database.GetConnection()
	if response.HandleError(res, err) {
		return
	}

	us := &services.User{DB: db}

	username := c.Param("username")

	user, err := us.Get(username)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		res.Fail().Status(http.StatusNotFound).Message("no user with username")
		return
	} else if response.HandleError(res, err) {
		return
	}

	res.Data(models.FromUser(user))
}

func getSelf(c *gin.Context) {
	res := response.From(c)

	t, exist := c.Get("user")
	user, ok := t.(*models.User)

	if !exist || !ok {
		res.Fail().Status(http.StatusNotFound).Message("no user with username")
		return
	}

	res.Data(models.FromUser(user))
}

func updateSelf(c *gin.Context) {
	res := response.From(c)

	db, err := database.GetConnection()
	if response.HandleError(res, err) {
		return
	}

	us := &services.User{DB: db}

	t, exist := c.Get("user")
	user, ok := t.(*models.User)

	if !exist || !ok {
		res.Fail().Status(http.StatusNotFound).Message("no user with username")
		return
	}

	var viewModel models.UserView
	if err = c.BindJSON(&viewModel); response.HandleError(res, err) {
		return
	}

	if err = viewModel.Valid(true); response.HandleError(res, err) {
		return
	}

	if viewModel.Password == "" {
		return
	}

	if us.IsValidCredentials(user, viewModel.Password) {
		response.HandleError(res, pufferpanel.ErrInvalidCredentials)
		return
	}

	viewModel.CopyToModel(user)

	if err = us.Update(user); response.HandleError(res, err) {
		return
	}

	res.Data(models.FromUser(user))
}

func updateUser(c *gin.Context) {
	res := response.From(c)

	db, err := database.GetConnection()
	if response.HandleError(res, err) {
		return
	}

	us := &services.User{DB: db}

	username := c.Param("username")

	var viewModel models.UserView
	if err = c.BindJSON(&viewModel); response.HandleError(res, err) {
		return
	}

	if err = viewModel.Valid(true); response.HandleError(res, err) {
		return
	}

	user, err := us.Get(username)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		res.Fail().Status(http.StatusNotFound).Message("no user with username")
		return
	} else if response.HandleError(res, err) {
		return
	}

	viewModel.CopyToModel(user)

	if err = us.Update(user); response.HandleError(res, err) {
		return
	}

	res.Data(models.FromUser(user))
}

func deleteUser(c *gin.Context) {
	res := response.From(c)

	db, err := database.GetConnection()
	if response.HandleError(res, err) {
		return
	}

	us := &services.User{DB: db}

	username := c.Param("username")

	user, err := us.Get(username)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		res.Fail().Status(http.StatusNotFound).Message("no user with username")
		return
	} else if response.HandleError(res, err) {
		return
	}

	if err = us.Delete(user.Username); response.HandleError(res, err) {
		return
	}

	res.Data(models.FromUser(user))
}

type UserSearch struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	PageLimit int    `json:"limit"`
	Page      int    `json:"page"`
}

func newUserSearch() *UserSearch {
	return &UserSearch{
		Username:  "*",
		Email:     "*",
		PageLimit: DefaultPageSize,
		Page:      1,
	}
}
