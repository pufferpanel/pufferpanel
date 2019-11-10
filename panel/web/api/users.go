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
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/panel/models"
	"github.com/pufferpanel/pufferpanel/v2/panel/services"
	"github.com/pufferpanel/pufferpanel/v2/panel/web/handlers"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"github.com/pufferpanel/pufferpanel/v2/scope"
	"net/http"
)

func registerUsers(g *gin.RouterGroup) {
	g.Handle("GET", "", handlers.OAuth2Handler(scope.UsersView, false), searchUsers)
	g.Handle("OPTIONS", "", response.CreateOptions("GET"))

	g.Handle("PUT", "/:username", handlers.OAuth2Handler(scope.UsersEdit, false), createUser)
	g.Handle("GET", "/:username", handlers.OAuth2Handler(scope.UsersView, false), getUser)
	g.Handle("POST", "/:username", handlers.OAuth2Handler(scope.UsersEdit, false), updateUser)
	g.Handle("DELETE", "/:username", handlers.OAuth2Handler(scope.UsersEdit, false), deleteUser)
	g.Handle("OPTIONS", "/:username", response.CreateOptions("PUT", "GET", "POST", "DELETE"))

	g.Handle("GET", "/:username/perms", handlers.OAuth2Handler(scope.UsersView, false), response.NotImplemented)
	g.Handle("PUT", "/:username/perms", handlers.OAuth2Handler(scope.UsersEdit, false), response.NotImplemented)
	g.Handle("OPTIONS", "/:username/perms", response.CreateOptions("PUT", "GET"))
}

func searchUsers(c *gin.Context) {
	var err error
	db := handlers.GetDatabase(c)
	us := &services.User{DB: db}

	search := newUserSearch()
	err = c.BindJSON(search)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	if search.PageLimit > MaxPageSize {
		search.PageLimit = MaxPageSize
	}

	var results *models.Users
	var total uint
	if results, total, err = us.Search(search.Username, search.Email, search.PageLimit, search.Page); response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, &UserSearchResponse{
		Users: models.FromUsers(results),
		Metadata: &response.Metadata{Paging: &response.Paging{
			Page:    search.Page,
			Size:    search.PageLimit,
			MaxSize: MaxPageSize,
			Total:   total,
		}},
	})
}

func createUser(c *gin.Context) {
	var err error
	db := handlers.GetDatabase(c)
	us := &services.User{DB: db}

	var viewModel models.UserView
	if err = c.BindJSON(&viewModel); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}
	viewModel.Username = c.Param("username")

	if err = viewModel.Valid(false); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	if viewModel.Password == "" {
		response.HandleError(c, pufferpanel.ErrFieldRequired("password"), http.StatusBadRequest)
		return
	}

	user := &models.User{}
	viewModel.CopyToModel(user)

	if err = us.Create(user); response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

func getUser(c *gin.Context) {
	db := handlers.GetDatabase(c)
	us := &services.User{DB: db}

	username := c.Param("username")

	user, err := us.Get(username)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, models.FromUser(user))
}

func updateUser(c *gin.Context) {
	db := handlers.GetDatabase(c)
	us := &services.User{DB: db}

	username := c.Param("username")

	var viewModel models.UserView
	if err := c.BindJSON(&viewModel); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	if err := viewModel.Valid(true); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	user, err := us.Get(username)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	viewModel.CopyToModel(user)

	if err = us.Update(user); response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

func deleteUser(c *gin.Context) {
	db := handlers.GetDatabase(c)
	us := &services.User{DB: db}

	username := c.Param("username")

	user, err := us.Get(username)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	if err = us.Delete(user.Username); response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

type UserSearch struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	PageLimit uint   `json:"limit"`
	Page      uint   `json:"page"`
}

type UserSearchResponse struct {
	Users []*models.UserView `json:"users"`
	*response.Metadata
}

func newUserSearch() *UserSearch {
	return &UserSearch{
		Username:  "*",
		Email:     "*",
		PageLimit: DefaultPageSize,
		Page:      1,
	}
}
