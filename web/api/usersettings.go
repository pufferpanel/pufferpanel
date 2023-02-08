package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/middleware/handlers"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"net/http"
)

func registerUserSettings(g *gin.RouterGroup) {
	g.Handle("GET", "", handlers.OAuth2Handler(pufferpanel.ScopeNone, false), getUserSettings)
	g.Handle("PUT", "/:key", handlers.OAuth2Handler(pufferpanel.ScopeNone, false), setUserSetting)
	g.Handle("OPTIONS", "", response.CreateOptions("GET", "PUT"))
}

// @Summary Get a user setting
// @Description Gets all settings specific to the current user
// @Produce json
// @Success 200 {object} models.UserSettingsView
// @Router /api/userSettings [get]
func getUserSettings(c *gin.Context) {
	db := middleware.GetDatabase(c)
	uss := &services.UserSettings{DB: db}

	t, exists := c.Get("user")
	user, ok := t.(*models.User)

	if !exists || !ok {
		response.HandleError(c, pufferpanel.ErrUnknownError, http.StatusInternalServerError)
		return
	}

	results, err := uss.GetAllForUser(user.ID)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, results)
}

// @Summary Update a user setting
// @Description Updates the value of a user setting
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty
// @Failure 400 {object} response.Error
// @Failure 500 {object} response.Error
// @Param key path string true "The config key"
// @Param value body models.ChangeUserSetting true "The new value for the setting"
// @Router /api/userSettings/{key} [PUT]
func setUserSetting(c *gin.Context) {
	key := c.Param("key")
	db := middleware.GetDatabase(c)
	uss := &services.UserSettings{DB: db}

	t, exists := c.Get("user")
	user, ok := t.(*models.User)

	if !exists || !ok {
		response.HandleError(c, pufferpanel.ErrUnknownError, http.StatusInternalServerError)
		return
	}

	var model models.ChangeUserSetting
	if err := c.BindJSON(&model); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	err := uss.Update(&models.UserSetting{
		Key:    key,
		UserID: user.ID,
		Value:  model.Value,
	})

	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}
