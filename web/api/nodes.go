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
	"github.com/pufferpanel/apufferi/response"
	"github.com/pufferpanel/apufferi/scope"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/web/handlers"
	"net/http"
	"strconv"
)

func registerNodes(g *gin.RouterGroup) {
	g.Handle("GET", "", handlers.OAuth2(scope.NodesView, false), getAllNodes)
	g.Handle("OPTIONS", "", response.CreateOptions("GET"))

	g.Handle("PUT", "/:id", handlers.OAuth2(scope.NodesEdit, false), createNode)
	g.Handle("GET", "/:id", handlers.OAuth2(scope.NodesView, false), getNode)
	g.Handle("POST", "/:id", handlers.OAuth2(scope.NodesEdit, false), updateNode)
	g.Handle("DELETE", "/:id", handlers.OAuth2(scope.NodesEdit, false), deleteNode)
	g.Handle("OPTIONS", "/:id", response.CreateOptions("PUT", "GET", "POST", "DELETE"))

	g.Handle("GET", "/:id/deployment", handlers.OAuth2(scope.NodesDeploy, false), response.NotImplemented)
	g.Handle("OPTIONS", "/:id/deployment", response.CreateOptions("GET"))

	g.Handle("POST", "/:id/reset", handlers.OAuth2(scope.NodesDeploy, false), response.NotImplemented)
	g.Handle("OPTIONS", "/:id/reset", response.CreateOptions("POST"))
}

func getAllNodes(c *gin.Context) {
	var err error
	res := response.From(c)

	db, err := database.GetConnection()
	if response.HandleError(res, err) {
		return
	}

	ns := &services.Node{DB: db}

	var nodes *models.Nodes
	if nodes, err = ns.GetAll(); response.HandleError(res, err) {
		return
	}

	data := models.FromNodes(nodes)

	res.Data(data)
}

func getNode(c *gin.Context) {
	var err error
	res := response.From(c)

	db, err := database.GetConnection()
	if response.HandleError(res, err) {
		return
	}

	ns := &services.Node{DB: db}

	id, ok := validateId(c, res)
	if !ok {
		return
	}

	node, exists, err := ns.Get(id)
	if response.HandleError(res, err) {
		return
	} else if !exists {
		res.Fail().Status(http.StatusNotFound)
		return
	}

	data := models.FromNode(node)

	res.Data(data)
}

func createNode(c *gin.Context) {
	var err error
	res := response.From(c)

	db, err := database.GetConnection()
	if response.HandleError(res, err) {
		return
	}

	ns := &services.Node{DB: db}

	model := &models.NodeView{}
	if err = c.BindJSON(model); response.HandleError(res, err) {
		return
	}

	if err = model.Valid(false); response.HandleError(res, err) {
		return
	}

	create := &models.Node{}
	model.CopyToModel(create)
	if err = ns.Create(create); response.HandleError(res, err) {
		return
	}

	res.Data(create)
}

func updateNode(c *gin.Context) {
	var err error
	res := response.From(c)

	db, err := database.GetConnection()
	if response.HandleError(res, err) {
		return
	}

	ns := &services.Node{DB: db}

	viewModel := &models.NodeView{}
	if err = c.BindJSON(viewModel); response.HandleError(res, err) {
		return
	}

	id, ok := validateId(c, res)
	if !ok {
		return
	}

	if err = viewModel.Valid(true); response.HandleError(res, err) {
		return
	}

	node, exists, err := ns.Get(id)
	if response.HandleError(res, err) {
		return
	} else if !exists {
		res.Fail().Status(http.StatusNotFound)
		return
	}

	viewModel.CopyToModel(node)
	if err = ns.Update(node); response.HandleError(res, err) {
		return
	}

	res.Data(node)
}

func deleteNode(c *gin.Context) {
	var err error
	res := response.From(c)

	db, err := database.GetConnection()
	if response.HandleError(res, err) {
		return
	}

	ns := &services.Node{DB: db}

	id, ok := validateId(c, res)
	if !ok {
		return
	}

	node, exists, err := ns.Get(id)
	if response.HandleError(res, err) {
		return
	} else if !exists {
		res.Fail().Status(http.StatusNotFound)
		return
	}

	err = ns.Delete(node.ID)
	if response.HandleError(res, err) {
		return
	}

	res.Data(node)
}

func validateId(c *gin.Context, response response.Builder) (uint, bool) {
	param := c.Param("id")

	id, err := strconv.Atoi(param)

	if err != nil || id <= 0 {
		response.Fail().Status(http.StatusBadRequest).Message("id must be a positive number")
		return 0, false
	}

	return uint(id), true
}
