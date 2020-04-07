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

package daemon

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"net/http"
)

func RegisterDaemonRoutes(e *gin.RouterGroup) {
	e.GET("", getStatusGET)
	e.HEAD("", getStatusHEAD)
	e.Handle("OPTIONS", "", response.CreateOptions("GET", "HEAD"))

	RegisterServerRoutes(e)
}

// Root godoc
// @Summary Is daemon up
// @Description Easy way to tell if the daemon is running is by using this endpoint
// @Accept json
// @Produce json
// @Success 200 {object} daemon.PufferdRunning "Service running"
// @Router / [get]
func getStatusGET(c *gin.Context) {
	c.JSON(http.StatusOK, &pufferpanel.PufferdRunning{Message: "pufferd is running"})
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
