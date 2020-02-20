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

package web

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v2/middleware"
	"github.com/pufferpanel/pufferpanel/v2/panel/web/api"
	"github.com/pufferpanel/pufferpanel/v2/panel/web/auth"
	"github.com/pufferpanel/pufferpanel/v2/panel/web/daemon"
	"github.com/pufferpanel/pufferpanel/v2/panel/web/handlers"
	"github.com/pufferpanel/pufferpanel/v2/panel/web/oauth2"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

var ClientPath string
var IndexFile string

var noHandle404 = []string{"/api/", "/oauth2/", "/daemon/"}

func RegisterRoutes(e *gin.Engine) {
	e.Use(func(c *gin.Context) {
		middleware.Recover(c)
	})

	ClientPath = viper.GetString("panel.web.files")
	IndexFile = ClientPath + "/index.html"

	api.RegisterRoutes(e.Group("/api", handlers.HasOAuth2Token))
	e.GET("/api/config", config)
	oauth2.RegisterRoutes(e.Group("/oauth2"))
	auth.RegisterRoutes(e.Group("/auth"))
	daemon.RegisterRoutes(e.Group("/daemon", handlers.HasOAuth2Token))

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
	//e.StaticFile("/", IndexFile)
	e.NoRoute( /*handlers.AuthMiddleware,*/ handle404)
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

	if strings.HasSuffix(c.Request.URL.Path, ".css") {
		c.Writer.Header().Set("Content-Type", "application/css")
		c.File(ClientPath + c.Request.URL.Path)
		return
	}

	c.File(IndexFile)
}

func config(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"branding": map[string]interface{}{
			"name": viper.GetString("panel.settings.companyName"),
		},
	})
}

func setContentType(contentType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", contentType)
		c.Next()
	}
}
