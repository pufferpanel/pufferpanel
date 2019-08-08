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
	builder "github.com/pufferpanel/apufferi/response"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/pufferpanel/pufferpanel/models/view"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/shared"
	"github.com/pufferpanel/pufferpanel/web/handlers"
	"net/http"
	"strconv"
)

func registerNodes(g *gin.RouterGroup) {
	g.Handle("GET", "", handlers.OAuth2("nodes.view", false), getAllNodes)
	g.Handle("OPTIONS", "", shared.CreateOptions("GET"))

	g.Handle("PUT", "/:id", handlers.OAuth2("nodes.edit", false), createNode)
	g.Handle("GET", "/:id", handlers.OAuth2("nodes.view", false), getNode)
	g.Handle("POST", "/:id", handlers.OAuth2("nodes.edit", false), updateNode)
	g.Handle("DELETE", "/:id", handlers.OAuth2("nodes.edit", false), deleteNode)
	g.Handle("OPTIONS", "/:id", shared.CreateOptions("PUT", "GET", "POST", "DELETE"))

	g.Handle("GET", "/:id/deployment", handlers.OAuth2("nodes.deploy", false), shared.NotImplemented)
	g.Handle("OPTIONS", "/:id/deployment", shared.CreateOptions("GET"))

	g.Handle("POST", "/:id/reset", handlers.OAuth2("nodes.deploy", false), shared.NotImplemented)
	g.Handle("OPTIONS", "/:id/reset", shared.CreateOptions("POST"))
}

func getAllNodes(c *gin.Context) {
	var err error
	response := builder.From(c)

	db, err := database.GetConnection()
	if shared.HandleError(response, err) {
		return
	}

	ns := &services.Node{DB: db}

	var nodes *models.Nodes
	if nodes, err = ns.GetAll(); shared.HandleError(response, err) {
		return
	}

	data := view.FromNodes(nodes)

	response.Data(data)
}

func getNode(c *gin.Context) {
	var err error
	response := builder.From(c)

	db, err := database.GetConnection()
	if shared.HandleError(response, err) {
		return
	}

	ns := &services.Node{DB: db}

	id, ok := validateId(c, response)
	if !ok {
		return
	}

	node, exists, err := ns.Get(id)
	if shared.HandleError(response, err) {
		return
	} else if !exists {
		response.Fail().Status(http.StatusNotFound).Send()
		return
	}

	data := view.FromNode(node)

	response.Data(data)
}

func createNode(c *gin.Context) {
	var err error
	response := builder.From(c)

	db, err := database.GetConnection()
	if shared.HandleError(response, err) {
		return
	}

	ns := &services.Node{DB: db}

	model := view.NodeViewModel{}
	if err = c.BindJSON(&model); shared.HandleError(response, err) {
		return
	}

	if err = model.Valid(false); shared.HandleError(response, err) {
		return
	}

	create := &models.Node{}
	model.CopyToModel(create)
	if err = ns.Create(create); shared.HandleError(response, err) {
		return
	}

	response.Data(create)
}

func updateNode(c *gin.Context) {
	var err error
	response := builder.From(c)

	db, err := database.GetConnection()
	if shared.HandleError(response, err) {
		return
	}

	ns := &services.Node{DB: db}

	viewModel := &view.NodeViewModel{}
	if err = c.BindJSON(viewModel); shared.HandleError(response, err) {
		return
	}

	id, ok := validateId(c, response)
	if !ok {
		return
	}

	if err = viewModel.Valid(true); shared.HandleError(response, err) {
		return
	}

	node, exists, err := ns.Get(id)
	if shared.HandleError(response, err) {
		return
	} else if !exists {
		response.Fail().Status(http.StatusNotFound).Send()
		return
	}

	viewModel.CopyToModel(node)
	if err = ns.Update(node); shared.HandleError(response, err) {
		return
	}

	response.Data(node).Send()
}

func deleteNode(c *gin.Context) {
	var err error
	response := builder.From(c)

	db, err := database.GetConnection()
	if shared.HandleError(response, err) {
		return
	}

	ns := &services.Node{DB: db}

	id, ok := validateId(c, response)
	if !ok {
		return
	}

	node, exists, err := ns.Get(id)
	if shared.HandleError(response, err) {
		return
	} else if !exists {
		response.Fail().Status(http.StatusNotFound)
		return
	}

	err = ns.Delete(node.ID)
	if shared.HandleError(response, err) {
		return
	}

	response.Data(node)
}

func validateId(c *gin.Context, response builder.Builder) (uint, bool) {
	param := c.Param("id")

	id, err := strconv.Atoi(param)

	if err != nil || id <= 0 {
		response.Fail().Status(http.StatusBadRequest).Message("id must be a positive number").Send()
		return 0, false
	}

	return uint(id), true
}
