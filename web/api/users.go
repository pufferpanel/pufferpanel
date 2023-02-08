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
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/middleware/handlers"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"net/http"
)

func registerUsers(g *gin.RouterGroup) {
	g.Handle("GET", "", handlers.OAuth2Handler(pufferpanel.ScopeUsersView, false), searchUsers)
	g.Handle("POST", "", handlers.OAuth2Handler(pufferpanel.ScopeUsersEdit, false), createUser)
	g.Handle("OPTIONS", "", response.CreateOptions("GET", "POST"))

	g.Handle("GET", "/:id", handlers.OAuth2Handler(pufferpanel.ScopeUsersView, false), getUser)
	g.Handle("POST", "/:id", handlers.OAuth2Handler(pufferpanel.ScopeUsersEdit, false), updateUser)
	g.Handle("DELETE", "/:id", handlers.OAuth2Handler(pufferpanel.ScopeUsersEdit, false), deleteUser)
	g.Handle("OPTIONS", "/:id", response.CreateOptions("GET", "POST", "DELETE"))

	g.Handle("GET", "/:id/perms", handlers.OAuth2Handler(pufferpanel.ScopeUsersView, false), getUserPerms)
	g.Handle("PUT", "/:id/perms", handlers.OAuth2Handler(pufferpanel.ScopeUsersEdit, false), setUserPerms)
	g.Handle("OPTIONS", "/:id/perms", response.CreateOptions("PUT", "GET"))
}

// @Summary Get users
// @Description Gets users, and allowing for filtering of users. * is a wildcard that can be used for text inputs
// @Accept json
// @Produce json
// @Success 200 {object} models.UserSearchResponse
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param body body models.UserSearch true "Filters to search on"
// @Router /api/users [get]
func searchUsers(c *gin.Context) {
	var err error
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	search := newUserSearch()
	err = c.ShouldBind(search)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	if search.PageLimit > MaxPageSize {
		search.PageLimit = MaxPageSize
	}

	var results *models.Users
	var total int64
	if results, total, err = us.Search(search.Username, search.Email, search.PageLimit, search.Page); response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, &models.UserSearchResponse{
		Users: models.FromUsers(results),
		Metadata: &response.Metadata{Paging: &response.Paging{
			Page:    search.Page,
			Size:    search.PageLimit,
			MaxSize: MaxPageSize,
			Total:   total,
		}},
	})
}

// @Summary Create user
// @Accept json
// @Produce json
// @Success 200 {object} models.UserView
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param body body models.UserView true "New user information"
// @Router /api/users [post]
func createUser(c *gin.Context) {
	var err error
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	var viewModel models.UserView
	if err = c.BindJSON(&viewModel); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

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

	resultModel := models.FromUser(user)

	c.JSON(http.StatusOK, resultModel)
}

// @Summary Get a user
// @Accept json
// @Produce json
// @Success 200 {object} models.UserView
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path uint true "User ID"
// @Router /api/users/{id} [get]
func getUser(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	var err error
	var id uint
	if id, err = cast.ToUintE(c.Param("id")); err != nil {
		response.HandleError(c, err, http.StatusBadRequest)
		return
	}

	user, err := us.GetById(id)
	if err != nil && err == gorm.ErrRecordNotFound {
		c.AbortWithStatus(http.StatusNotFound)
		return
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, models.FromUser(user))
}

// @Summary Update user
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path uint true "User ID"
// @Param body body models.UserView true "New user information"
// @Router /api/users/{id} [post]
func updateUser(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	var err error
	var id uint
	if id, err = cast.ToUintE(c.Param("id")); err != nil {
		response.HandleError(c, err, http.StatusBadRequest)
		return
	}

	var viewModel models.UserView
	if err := c.BindJSON(&viewModel); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	if err := viewModel.Valid(true); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	user, err := us.GetById(id)
	if err != nil && err == gorm.ErrRecordNotFound {
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

// @Summary Delete user
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path uint true "User ID"
// @Router /api/users/{id} [delete]
func deleteUser(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	var err error
	var id uint
	if id, err = cast.ToUintE(c.Param("id")); err != nil {
		response.HandleError(c, err, http.StatusBadRequest)
		return
	}

	user, err := us.GetById(id)
	if err != nil && err == gorm.ErrRecordNotFound {
		c.AbortWithStatus(http.StatusNotFound)
		return
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	if err = us.Delete(user); response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Gets user permissions
// @Accept json
// @Produce json
// @Success 200 {object} models.PermissionView
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path uint true "User ID"
// @Router /api/users/{id}/perms [get]
func getUserPerms(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}
	ps := &services.Permission{DB: db}

	var err error
	var id uint
	if id, err = cast.ToUintE(c.Param("id")); err != nil {
		response.HandleError(c, err, http.StatusBadRequest)
		return
	}

	user, err := us.GetById(id)
	if err != nil && err == gorm.ErrRecordNotFound {
		c.AbortWithStatus(http.StatusNotFound)
		return
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	perms, err := ps.GetForUserAndServer(user.ID, nil)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, models.FromPermission(perms))
}

// @Summary Sets user permissions
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path uint true "User ID"
// @Param body body models.PermissionView true "New permissions"
// @Router /api/users/{id}/perms [put]
func setUserPerms(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}
	ps := &services.Permission{DB: db}

	var err error
	var id uint
	if id, err = cast.ToUintE(c.Param("id")); err != nil {
		response.HandleError(c, err, http.StatusBadRequest)
		return
	}

	viewModel := &models.PermissionView{}
	err = c.BindJSON(viewModel)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	user, err := us.GetById(id)
	if err != nil && err == gorm.ErrRecordNotFound {
		c.AbortWithStatus(http.StatusNotFound)
		return
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	perms, err := ps.GetForUserAndServer(user.ID, nil)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	viewModel.CopyTo(perms, true)

	err = ps.UpdatePermissions(perms)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

func newUserSearch() *models.UserSearch {
	return &models.UserSearch{
		Username:  "*",
		Email:     "*",
		PageLimit: DefaultPageSize,
		Page:      1,
	}
}
