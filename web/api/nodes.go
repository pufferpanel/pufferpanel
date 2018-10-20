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
	"github.com/pufferpanel/pufferpanel/models/view"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/shared"
	builder "github.com/pufferpanel/apufferi/http"
	"net/http"
	"strconv"
)

func registerNodes(g *gin.RouterGroup) {
	g.Handle("GET", "", GetAllNodes)
	g.Handle("OPTIONS", "", shared.CreateOptions("GET"))

	g.Handle("PUT", "/:id", shared.NotImplemented)
	g.Handle("GET", "/:id", GetNode)
	g.Handle("POST", "/:id", shared.NotImplemented)
	g.Handle("DELETE", "/:id", shared.NotImplemented)
	g.Handle("OPTIONS", "/:id", shared.CreateOptions("PUT", "GET", "POST", "DELETE"))
}

func GetAllNodes(c *gin.Context) {
	response := builder.Respond(c)

	ns, err := services.GetNodeService()
	if err != nil {
		response.Fail().Message(err.Error()).Send()
		return
	}

	nodes, err := ns.GetAll()
	if err != nil {
		response.Fail().Message(err.Error()).Send()
		return
	}

	data := view.FromNodes(nodes)

	response.Data(data).Send()
}

func GetNode(c *gin.Context) {
	response := builder.Respond(c)

	ns, err := services.GetNodeService()
	if err != nil {
		response.Fail().Message(err.Error()).Send()
		return
	}

	param := c.Param("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		response.Fail().Message("id must be a number").Send()
		return
	}

	node, exists, err := ns.Get(id)
	if err != nil {
		response.Fail().Message(err.Error()).Send()
		return
	}

	if !exists {
		response.Fail().Status(http.StatusNotFound).Message("no node with given id").Send()
		return
	}

	data := view.FromNode(node)

	response.Data(data).Send()
}