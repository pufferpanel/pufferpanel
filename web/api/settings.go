package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/config"
	"github.com/pufferpanel/pufferpanel/v2/middleware/handlers"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"net/http"
)

func registerSettings(g *gin.RouterGroup) {
	g.Handle("GET", "/:key", handlers.OAuth2Handler(pufferpanel.ScopeServersAdmin, false), getSetting)
	g.Handle("PUT", "/:key", handlers.OAuth2Handler(pufferpanel.ScopeServersAdmin, false), setSetting)
	g.Handle("OPTIONS", "", response.CreateOptions("GET", "PUT"))
}

// @Summary Get a panel setting
// @Description Gets the value currently being used for the specified config key
// @Produce json
// @Success 200 {object} models.SettingResponse
// @Param key path string true "The config key"
// @Router /api/settings/{key} [get]
func getSetting(c *gin.Context) {
	key := c.Param("key")

	response := &models.SettingResponse{
		Value: config.GetString(key),
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Update a panel setting
// @Description Updates the value of a panel setting
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty
// @Failure 400 {object} response.Error
// @Failure 500 {object} response.Error
// @Param key path string true "The config key"
// @Param value body models.ChangeSetting true "The new value for the setting"
// @Router /api/self [PUT]
func setSetting(c *gin.Context) {
	key := c.Param("key")

	var model models.ChangeSetting
	if err := c.BindJSON(&model); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	if err := config.Set(key, model.Value); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	c.Status(http.StatusNoContent)
}
