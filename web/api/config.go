package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v2/config"
	"io/ioutil"
	"net/http"
	"strings"
)

func panelConfig(c *gin.Context) {
	themes := []string{}
	files, err := ioutil.ReadDir(config.GetString("panel.web.files") + "/theme")
	if err != nil {
		themes = append(themes, "PufferPanel")
	} else {
		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".tar") {
				themes = append(themes, f.Name()[:len(f.Name())-4])
			}
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"themes": map[string]interface{}{
			"active": config.GetString("panel.settings.defaultTheme"),
			"available": themes,
		},
		"branding": map[string]interface{}{
			"name": config.GetString("panel.settings.companyName"),
		},
		"registrationEnabled": true,
	})
}