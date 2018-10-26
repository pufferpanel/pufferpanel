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
	builder "github.com/pufferpanel/apufferi/http"
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/pufferpanel/pufferpanel/models/view"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/shared"
	"net/http"
	"strconv"
)

func registerServers(g *gin.RouterGroup) {
	g.Handle("GET", "", SearchServers)
	g.Handle("OPTIONS", "", shared.CreateOptions("GET"))

	g.Handle("PUT", "", shared.NotImplemented)
	g.Handle("GET", "/:id", GetServer)
	g.Handle("POST", "/:id", shared.NotImplemented)
	g.Handle("DELETE", "/:id", shared.NotImplemented)
	g.Handle("OPTIONS", "/:id", shared.CreateOptions("PUT", "GET", "POST", "DELETE"))
}

func SearchServers (c *gin.Context) {
	var ss *services.ServerService
	var err error
	response := builder.Respond(c)

	nodeQuery := c.DefaultQuery("node", "0")
	nameFilter := c.DefaultQuery("name", "*")
	pageSizeQuery := c.DefaultQuery("limit", strconv.Itoa(DefaultPageSize))
	pageQuery := c.DefaultQuery("page", strconv.Itoa(1))

	pageSize, err := strconv.Atoi(pageSizeQuery)
	if err != nil || pageSize <= 0 {
		response.Fail().Status(http.StatusBadRequest).Message("page size must be a positive number").Send()
		return
	}

	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	page, err := strconv.Atoi(pageQuery)
	if err != nil || page <= 0 {
		response.Fail().Status(http.StatusBadRequest).Message("page must be a positive number").Send()
		return
	}

	node, err := strconv.Atoi(nodeQuery)
	if err != nil || page <= 0 {
		response.Fail().Status(http.StatusBadRequest).Message("node id is invalid").Send()
		return
	}

	if ss, err = services.GetServerService(); shared.HandleError(response, err) {
		return
	}

	var results *models.Servers
	if results, err = ss.Search(uint(node), nameFilter, uint(pageSize), uint(page)); shared.HandleError(response, err) {
		return
	}

	response.PageInfo(uint(page), uint(pageSize), MaxPageSize).Data(view.FromServers(results)).Send()
}

func GetServer(c *gin.Context) {
	var ss *services.ServerService
	var err error
	response := builder.Respond(c)

	id := c.Param("id")

	if ss, err = services.GetServerService(); shared.HandleError(response, err) {
		return
	}

	var result *models.Server
	var exists bool
	if result, exists, err = ss.Get(id); shared.HandleError(response, err) {
		return
	}

	if !exists {
		response.Fail().Status(http.StatusNotFound).Send()
		return
	}

	response.Data(result).Send()
}