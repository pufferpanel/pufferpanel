/*
 Copyright 2016 Padduck, LLC

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

package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v2/daemon"
	"github.com/pufferpanel/pufferpanel/v2/daemon/routing/server"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/middleware"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"net/http"
	"strings"
)

// @title Pufferd API
// @version 2.0
// @description PufferPanel daemon service
// @contact.name PufferPanel
// @contact.url https://pufferpanel.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func ConfigureWeb() *gin.Engine {
	r := gin.New()
	{
		r.Use(gin.Recovery())
		r.Use(gin.LoggerWithWriter(logging.Info().Writer()))
		r.Use(func(c *gin.Context) {
			if c.GetHeader("Connection") == "Upgrade" {
				return
			}
			if strings.HasPrefix(c.Request.URL.Path, "/swagger/") {
				return
			}
			middleware.ResponseAndRecover(c)
		})
		RegisterRoutes(r)
		server.RegisterRoutes(r.Group("/"))
	}

	return r
}

func RegisterRoutes(e *gin.Engine) {
	e.GET("", getStatusGET)
	e.HEAD("", getStatusHEAD)
	e.Handle("OPTIONS", "", response.CreateOptions("GET", "HEAD"))
}

// Root godoc
// @Summary Is daemon up
// @Description Easy way to tell if the daemon is running is by using this endpoint
// @Accept json
// @Produce json
// @Success 200 {object} daemon.PufferdRunning "Service running"
// @Router / [get]
func getStatusGET(c *gin.Context) {
	c.JSON(http.StatusOK, &daemon.PufferdRunning{Message: "pufferd is running"})
}

// Root godoc
// @Summary Is daemon up
// @Description Easy way to tell if the daemon is running is by using this endpoint
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty "Service running"
// @Router / [head]
func getStatusHEAD(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
