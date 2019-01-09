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
	"fmt"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferpanel/web/api"
	"github.com/pufferpanel/pufferpanel/web/auth"
	"github.com/pufferpanel/pufferpanel/web/oauth2"
	"github.com/pufferpanel/pufferpanel/web/server"
	"path/filepath"
	"strings"
)

func RegisterRoutes(e *gin.Engine) {
	e.HTMLRender = loadTemplates()

	e.Group("/assets").Static("", "assets/web")
	api.RegisterRoutes(e.Group("/api"))
	oauth2.RegisterRoutes(e.Group("/oauth2"))
	auth.RegisterRoutes(e.Group("/auth"))
	server.RegisterRoutes(e.Group("/server"))

	e.Handle("GET", "/", func (c *gin.Context) {
		c.Redirect(302, "/server")
	})
}

func loadTemplates() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	prefix := "assets/web/"

	layouts, err := filepath.Glob(prefix + "base.html")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(prefix + "**/*.html")
	if err != nil {
		panic(err.Error())
	}

	for _, include := range includes {
		templateName := strings.TrimPrefix(include, prefix)
		if templateName == include {
			templateName = strings.TrimPrefix(include, strings.Replace(prefix, "/", "\\", 2))
		}
		templateName = strings.Replace(strings.TrimSuffix(templateName, ".html"), "\\", "/", -1)
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		logging.Debugf(fmt.Sprintf("Adding template [%s] with %s", templateName, files))
		r.AddFromFiles(templateName, files...)
	}
	return r
}
