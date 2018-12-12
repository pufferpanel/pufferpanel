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
	"github.com/pufferpanel/pufferpanel/web/handlers"
	"github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

func registerServers(g *gin.RouterGroup) {
	g.Handle("GET", "", handlers.OAuth2("servers.view", false), searchServers)
	g.Handle("OPTIONS", "", shared.CreateOptions("GET"))

	g.Handle("POST", "", handlers.OAuth2("servers.edit", false), createServer)
	g.Handle("GET", "/:id", handlers.OAuth2("servers.view", true), getServer)
	g.Handle("PUT", "/:id", handlers.OAuth2("servers.edit", false), createServer)
	g.Handle("DELETE", "/:id", handlers.OAuth2("servers.edit", false), deleteServer)
	g.Handle("GET", "/:id/users", handlers.OAuth2("servers.edit", true), getServerUsers)
	g.Handle("POST", "/:id/users", handlers.OAuth2("servers.edit", true), editServerUsers)
	g.Handle("OPTIONS", "/:id", shared.CreateOptions("PUT", "GET", "POST", "DELETE"))
}

func searchServers(c *gin.Context) {
	var ss services.ServerService
	var err error
	response := builder.Respond(c)

	username := c.DefaultQuery("username", "")
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

	//see if user has access to view all others, otherwise we can't permit search without their username
	os, _ := services.GetOAuthService()
	if ci, global, _ := os.HasRights(c.GetString("accessToken"), nil, "servers.view"); !global {
		username = ci.User.Username
	}

	var results *models.Servers
	if results, err = ss.Search(username, uint(node), nameFilter, uint(pageSize), uint(page)); shared.HandleError(response, err) {
		return
	}

	response.PageInfo(uint(page), uint(pageSize), MaxPageSize).Data(view.FromServers(results)).Send()
}

func getServer(c *gin.Context) {
	var ss services.ServerService
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

func createServer(c *gin.Context) {
	var ss services.ServerService
	var ns services.NodeService
	var err error
	response := builder.Respond(c)

	serverId := c.Param("id")
	if serverId == "" {
		serverId = uuid.NewV4().String()[:8]
	}

	postBody := view.ServerViewModel{}
	err = c.Bind(&postBody)
	postBody.Identifier = serverId
	if err != nil {
		response.Status(http.StatusBadRequest).Message(err.Error()).Fail().Send()
		return
	}

	if ss, err = services.GetServerService(); shared.HandleError(response, err) {
		return
	}

	if ns, err = services.GetNodeService(); shared.HandleError(response, err) {
		return
	}

	node, exists, err := ns.Get(postBody.NodeId)

	if shared.HandleError(response, err) {
		return
	}

	if !exists {
		response.Status(http.StatusBadRequest).Message("no node with given id").Fail().Send()
	}

	server := &models.Server{}
	postBody.CopyToModel(server)

	server.NodeID = node.ID

	err = ss.Create(server, postBody.Data)
	if err != nil {
		response.Status(http.StatusInternalServerError).Message(err.Error()).Fail().Send()
		return
	}

	postBody.Data = nil
	response.Data(postBody).Send()
}

func deleteServer(c *gin.Context) {
	var ss services.ServerService
	var err error
	response := builder.Respond(c)

	serverId := c.Param("id")

	if ss, err = services.GetServerService(); shared.HandleError(response, err) {
		return
	}

	server, exists, err := ss.Get(serverId)
	if shared.HandleError(response, err) {
		return
	}

	if !exists {
		response.Status(http.StatusNotFound).Message("no server with given id").Fail().Send()
	}

	err = ss.Delete(server.ID)
	if shared.HandleError(response, err) {
		return
	} else {
		v := view.FromServer(server)
		response.Status(http.StatusOK).Data(v).Send()
	}
}

func getServerUsers(c *gin.Context){

}

func editServerUsers(c *gin.Context) {

}