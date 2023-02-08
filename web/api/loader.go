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
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/middleware/handlers"
)

const MaxPageSize = 100
const DefaultPageSize = 20

func RegisterRoutes(rg *gin.RouterGroup) {

	rg.Use(func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Next()
	})

	rg.Use(middleware.ResponseAndRecover)
	rg.Use(middleware.NeedsDatabase)
	registerNodes(rg.Group("/nodes", handlers.HasOAuth2Token))
	registerServers(rg.Group("/servers", handlers.HasOAuth2Token))
	registerUsers(rg.Group("/users", handlers.HasOAuth2Token))
	registerTemplates(rg.Group("/templates", handlers.HasOAuth2Token))
	registerSelf(rg.Group("/self", handlers.HasOAuth2Token))
	registerSettings(rg.Group("/settings", handlers.HasOAuth2Token))
	registerUserSettings(rg.Group("/userSettings", handlers.HasOAuth2Token))

	rg.GET("/config", panelConfig)
}
