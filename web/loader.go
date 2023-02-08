/*
 Copyright 2022 (c) PufferPanel

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

package web

import (
	"fmt"
	_ "github.com/alecthomas/template"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/web/api"
	"github.com/pufferpanel/pufferpanel/v3/web/auth"
	"github.com/pufferpanel/pufferpanel/v3/web/daemon"
	"github.com/pufferpanel/pufferpanel/v3/web/oauth2"
	_ "github.com/pufferpanel/pufferpanel/v3/web/swagger"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "github.com/swaggo/swag"
	"net/http"
	"strings"
)

var ClientPath string
var IndexFile string

var noHandle404 = []string{"/api/", "/oauth2/", "/daemon/", "/proxy/"}

// @title PufferPanel API
// @version 3.0
// @description PufferPanel API interface for both the panel and daemon. Endpoints starting with /daemon or /proxy are for nodes.
// @contact.name PufferPanel
// @contact.url https://pufferpanel.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func RegisterRoutes(e *gin.Engine) {
	e.Use(func(c *gin.Context) {
		middleware.Recover(c)
	})

	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if config.DaemonEnabled.Value() {
		daemon.RegisterDaemonRoutes(e.Group("/daemon"))
	}

	if config.PanelEnabled.Value() {
		ClientPath = config.WebRoot.Value()
		IndexFile = ClientPath + "/index.html"

		api.RegisterRoutes(e.Group("/api"))
		e.GET("/manifest.json", webManifest)
		oauth2.RegisterRoutes(e.Group("/oauth2"))
		auth.RegisterRoutes(e.Group("/auth"))

		css := e.Group("/css")
		{
			css.Use(gzip.Gzip(gzip.DefaultCompression))
			css.StaticFS("", http.Dir(ClientPath+"/css"))
		}
		fonts := e.Group("/fonts")
		{
			fonts.Use(gzip.Gzip(gzip.DefaultCompression))
			fonts.StaticFS("", http.Dir(ClientPath+"/fonts"))
		}
		img := e.Group("/img")
		{
			img.StaticFS("", http.Dir(ClientPath+"/img"))
		}
		js := e.Group("/js", setContentType("application/javascript"))
		{
			js.Use(gzip.Gzip(gzip.DefaultCompression))
			js.StaticFS("", http.Dir(ClientPath+"/js"))
		}
		e.StaticFile("/favicon.png", ClientPath+"/favicon.png")
		e.StaticFile("/favicon.ico", ClientPath+"/favicon.ico")
		e.NoRoute(handle404)
	}
}

func handle404(c *gin.Context) {
	for _, v := range noHandle404 {
		if strings.HasPrefix(c.Request.URL.Path, v) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	}

	if strings.HasSuffix(c.Request.URL.Path, ".js") {
		c.Writer.Header().Set("Content-Type", "application/javascript")
		c.File(ClientPath + c.Request.URL.Path)
		return
	}

	if strings.HasSuffix(c.Request.URL.Path, ".json") {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.File(ClientPath + c.Request.URL.Path)
		return
	}

	if strings.HasSuffix(c.Request.URL.Path, ".css") {
		c.Writer.Header().Set("Content-Type", "text/css")
		c.File(ClientPath + c.Request.URL.Path)
		return
	}

	if strings.HasSuffix(c.Request.URL.Path, ".tar") {
		c.Writer.Header().Set("Content-Type", "application/x-tar")
		c.File(ClientPath + c.Request.URL.Path)
		return
	}

	c.File(IndexFile)
}

func webManifest(c *gin.Context) {
	iconSizes := []int{72, 96, 128, 144, 152, 192, 384, 512}
	icons := make([]map[string]interface{}, len(iconSizes))

	for i, s := range iconSizes {
		icons[i] = map[string]interface{}{
			"src":   fmt.Sprintf("img/appicons/%d.png", s),
			"sizes": fmt.Sprintf("%dx%d", s, s),
			"type":  "image/png",
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"name":             config.CompanyName.Value(),
		"short_name":       config.CompanyName.Value(),
		"background_color": "#fff",
		"display":          "standalone",
		"scope":            "/",
		"start_url":        "/servers",
		"icons":            icons,
	})
}

func setContentType(contentType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", contentType)
	}
}
