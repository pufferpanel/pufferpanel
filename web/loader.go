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
	"github.com/gin-gonic/gin/binding"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/client/frontend/dist"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/web/api"
	"github.com/pufferpanel/pufferpanel/v3/web/auth"
	"github.com/pufferpanel/pufferpanel/v3/web/daemon"
	"github.com/pufferpanel/pufferpanel/v3/web/oauth2"
	_ "github.com/pufferpanel/pufferpanel/v3/web/swagger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/swaggo/swag"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

var noHandle404 = []string{"/api/", "/oauth2/", "/daemon/", "/proxy/"}
var clientFiles fs.ReadFileFS

// RegisterRoutes Registers all routes
// @title PufferPanel API
// @version 3.0
// @description PufferPanel API interface for both the panel and daemon.
// @contact.name PufferPanel
// @contact.url https://pufferpanel.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @Accept json
// @Produce json
// @description.markdown
// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl /oauth2/token
// @scope.none No scope needed
// @scope.admin Administrator, full rights to all actions
// @scope.login Allows logging into the panel
// @scope.oauth2.auth Scope to validate another OAuth2 credential
// @scope.nodes.view Allows viewing nodes
// @scope.nodes.edit Allows editing of node connection information
// @scope.nodes.deploy Allows getting the config of a node for deployment
// @scope.self.edit Allows editing of personal account
// @scope.self.clients Allows creating OAuth2 clients under the account
// @scope.server.admin Admin access to a server (full permissions)
// @scope.server.view Allows viewing a server
// @scope.server.create Allows creating servers
// @scope.server.delete Allows deleting servers
// @scope.server.edit.definition Allows editing a server's definition
// @scope.server.edit.data Allows editing the values of variables
// @scope.server.edit.flags Allows changing flags on the server
// @scope.server.edit.name Allows editing of a server name
// @scope.server.view.definition Allows viewing a server's definition
// @scope.server.view.data Allows viewing a server's variables
// @scope.server.view.flags Allows viewing a server's flags
// @scope.server.clients.view Allows viewing OAuth2 clients associated to a server
// @scope.server.clients.edit Allows editing OAuth2 clients associated to a server
// @scope.server.clients.add Allows adding a new OAuth2 client to a server
// @scope.server.clients.delete Allows deleting OAuth2 clients associated to a server
// @scope.server.users.view Allows viewing users associated to a server
// @scope.server.users.edit Allows editing user permissions to a server
// @scope.server.users.add Allows adding a new user to a server
// @scope.server.users.delete Allows removing users from to a server
// @scope.server.tasks.view Allows viewing tasks associated to a server
// @scope.server.tasks.edit Allows editing tasks associated to a server
// @scope.server.tasks.add Allows adding a new tasks to a server
// @scope.server.tasks.delete Allows deleting tasks from to a server
// @scope.server.reload Allows reloading of a server's definition from disk
// @scope.server.start Allow starting a server
// @scope.server.stop Allows stopping a server
// @scope.server.kill Allows killing a server
// @scope.server.install Allows using the "Install" button for a server
// @scope.server.files.get Allows viewing and downloading files for a server through the File Manager
// @scope.server.files.edit Allows editing files for a server through the File Manager
// @scope.server.sftp Allows connection to a server over SFTP
// @scope.server.console Allows viewing the console of a server
// @scope.server.console.send Allows sending commands to a server's console
// @scope.server.stats Allows getting stats of a server like CPU and memory usage
// @scope.server.status Allows getting the status of a server
// @scope.panel.settings.edit Allows for viewing and editing of panel settings
// @scope.templates.view Allows viewing templates
// @scope.templates.local.edit Allows editing of templates in the local repo
// @scope.templates.repo.add Allows adding a new template repo
// @scope.templates.repo.delete Allows deleting of a template repo
// @scope.users.view Allows viewing all registered users
// @scope.users.edit Allows editing of all users
func RegisterRoutes(e *gin.Engine) {
	e.Use(func(c *gin.Context) {
		middleware.Recover(c)
	})

	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.DefaultModelsExpandDepth(0), ginSwagger.DeepLinking(false)))

	if config.DaemonEnabled.Value() {
		daemon.RegisterDaemonRoutes(e.Group("/daemon"))
	}

	if config.PanelEnabled.Value() {
		api.RegisterRoutes(e.Group("/api"))
		e.GET("/manifest.json", webManifest)
		oauth2.RegisterRoutes(e.Group("/oauth2"))
		auth.RegisterRoutes(e.Group("/auth"))

		clientFiles = dist.ClientFiles
		if config.WebRoot.Value() != "" {
			clientFiles = pufferpanel.NewMergedFS(os.DirFS(config.WebRoot.Value()), dist.ClientFiles)
		}

		css := e.Group("/css")
		{
			css.Use(gzip.Gzip(gzip.DefaultCompression))
			css.Use(setContentType("text/css"))
			f, err := fs.Sub(clientFiles, "css")
			if err != nil {
				panic(err)
			}
			css.StaticFS("", http.FS(f))
		}
		fonts := e.Group("/fonts")
		{
			fonts.Use(gzip.Gzip(gzip.DefaultCompression))
			f, err := fs.Sub(clientFiles, "fonts")
			if err != nil {
				panic(err)
			}
			fonts.StaticFS("", http.FS(f))
		}
		img := e.Group("/img")
		{
			f, err := fs.Sub(clientFiles, "img")
			if err != nil {
				panic(err)
			}
			img.StaticFS("", http.FS(f))
		}
		js := e.Group("/js")
		{
			js.Use(gzip.Gzip(gzip.DefaultCompression))
			js.Use(setContentType("application/javascript"))
			f, err := fs.Sub(clientFiles, "js")
			if err != nil {
				panic(err)
			}
			js.StaticFS("", http.FS(f))
		}
		wasm := e.Group("/wasm")
		{
			wasm.Use(gzip.Gzip(gzip.DefaultCompression))
			wasm.Use(setContentType("application/wasm"))
			f, err := fs.Sub(clientFiles, "wasm")
			if err != nil {
				panic(err)
			}
			wasm.StaticFS("", http.FS(f))
		}
		theme := e.Group("/theme")
		{
			theme.Use(setContentType("application/x-tar"))
			f, err := fs.Sub(clientFiles, "theme")
			if err != nil {
				panic(err)
			}
			theme.StaticFS("", http.FS(f))
		}
		e.StaticFileFS("/favicon.png", "favicon.png", http.FS(clientFiles))
		e.StaticFileFS("/favicon.ico", "favicon.ico", http.FS(clientFiles))
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

	file, err := clientFiles.ReadFile("index.html")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Data(http.StatusOK, binding.MIMEHTML, file)
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
