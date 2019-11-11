package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/panel/models"
	"github.com/pufferpanel/pufferpanel/v2/panel/services"
	"github.com/pufferpanel/pufferpanel/v2/panel/web/handlers"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"github.com/pufferpanel/pufferpanel/v2/scope"
	"net/http"
)

func registerSelf(g *gin.RouterGroup) {
	g.Handle("GET", "", handlers.OAuth2Handler(scope.Login, false), getSelf)
	g.Handle("PUT", "", handlers.OAuth2Handler(scope.Login, false), updateSelf)
}

func getSelf(c *gin.Context) {
	t, exist := c.Get("user")
	user, ok := t.(*models.User)

	if !exist || !ok {
		response.HandleError(c, pufferpanel.ErrUnknownError, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, models.FromUser(user))
}

func updateSelf(c *gin.Context) {
	db := handlers.GetDatabase(c)
	us := &services.User{DB: db}

	t, exist := c.Get("user")
	user, ok := t.(*models.User)

	if !exist || !ok {
		response.HandleError(c, pufferpanel.ErrUnknownError, http.StatusInternalServerError)
		return
	}

	var viewModel models.UserView
	if err := c.BindJSON(&viewModel); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	if err := viewModel.Valid(true); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	if viewModel.Password == "" {
		response.HandleError(c, pufferpanel.ErrFieldRequired("password"), http.StatusBadRequest)
		return
	}

	if !us.IsValidCredentials(user, viewModel.Password) {
		response.HandleError(c, pufferpanel.ErrInvalidCredentials, http.StatusInternalServerError)
		return
	}

	viewModel.CopyToModel(user)

	if viewModel.NewPassword != "" {
		err := user.SetPassword(viewModel.NewPassword)
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}
	}

	if err := us.Update(user); response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}
