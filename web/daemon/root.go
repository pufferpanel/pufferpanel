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

package daemon

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"net/http"
	"time"
)

func RegisterDaemonRoutes(e *gin.RouterGroup) {
	e.GET("", getStatusGET)
	e.HEAD("", getStatusHEAD)
	e.Handle("OPTIONS", "", response.CreateOptions("GET", "HEAD"))
	e.GET("features", getFeatures)

	RegisterServerRoutes(e)
}

func getStatusGET(c *gin.Context) {
	c.JSON(http.StatusOK, &pufferpanel.DaemonRunning{Message: "daemon is running"})
}

func getStatusHEAD(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func getFeatures(c *gin.Context) {
	features := []string{}

	if testDocker() {
		features = append(features, "docker")
	}

	c.JSON(http.StatusOK, Features{Features: features})
}

func testDocker() bool {
	d, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = d.Ping(ctx)
	if err != nil {
		return false
	}
	return true
}

type Features struct {
	Features []string `json:"features"`
}
