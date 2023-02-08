package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/spf13/cast"
	"net/http"
)

func registerSettings(g *gin.RouterGroup) {
	g.Handle("GET", "/:key", middleware.RequiresPermission(pufferpanel.ScopeSettings, false), getSetting)
	g.Handle("PUT", "/:key", middleware.RequiresPermission(pufferpanel.ScopeSettings, false), setSetting)
	g.Handle("POST", "", middleware.RequiresPermission(pufferpanel.ScopeSettings, false), setSettings)
	g.Handle("OPTIONS", "", response.CreateOptions("GET", "PUT"))
}

// @Summary Value a panel setting
// @Description Gets the value currently being used for the specified config key
// @Produce json
// @Success 200 {object} models.SettingResponse
// @Param key path string true "The config key"
// @Router /api/settings/{key} [get]
func getSetting(c *gin.Context) {
	key := c.Param("key")

	for _, v := range editableStringEntries {
		if v.Key() == key {
			c.JSON(http.StatusOK, models.SettingResponse{Value: v.Value()})
			return
		}
	}

	for _, v := range editableBoolEntries {
		if v.Key() == key {
			c.JSON(http.StatusOK, models.SettingResponse{Value: v.Value()})
			return
		}
	}

	for _, v := range editableIntEntries {
		if v.Key() == key {
			c.JSON(http.StatusOK, models.SettingResponse{Value: v.Value()})
			return
		}
	}

	c.Status(http.StatusNoContent)
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
	var err error
	if err = c.BindJSON(&model); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	for _, v := range editableStringEntries {
		if v.Key() == key {
			err = v.Set(cast.ToString(model.Value), true)
		}
	}

	for _, v := range editableBoolEntries {
		if v.Key() == key {
			err = v.Set(cast.ToBool(model.Value), true)
		}
	}

	for _, v := range editableIntEntries {
		if v.Key() == key {
			err = v.Set(cast.ToInt(model.Value), true)
		}
	}

	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

func setSettings(c *gin.Context) {
	var settings map[string]interface{}
	var err error
	if err = c.BindJSON(&settings); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	for key, value := range settings {
		for _, v := range editableStringEntries {
			if v.Key() == key {
				err = v.Set(cast.ToString(value), true)
			}
		}

		for _, v := range editableBoolEntries {
			if v.Key() == key {
				err = v.Set(cast.ToBool(value), true)
			}
		}

		for _, v := range editableIntEntries {
			if v.Key() == key {
				err = v.Set(cast.ToInt(value), true)
			}
		}
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}
	}

	c.Status(http.StatusNoContent)
}

var editableStringEntries = []config.StringEntry{
	config.EmailDomain,
	config.EmailFrom,
	config.EmailHost,
	config.EmailKey,
	config.EmailPassword,
	config.EmailProvider,
	config.EmailUsername,
	config.CompanyName,
	config.DefaultTheme,
	config.ThemeSettings,
	config.MasterUrl,
}
var editableBoolEntries = []config.BoolEntry{
	config.RegistrationEnabled,
}
var editableIntEntries = []config.IntEntry{}
