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
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/gofrs/uuid/v5"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"github.com/pufferpanel/pufferpanel/v3/web/daemon"
	"net/http"
	"strconv"
	"strings"
)

func registerNodes(g *gin.RouterGroup) {
	g.Handle("GET", "", middleware.RequiresPermission(pufferpanel.ScopeNodesView, false), getAllNodes)
	g.Handle("POST", "", middleware.RequiresPermission(pufferpanel.ScopeNodesCreate, false), createNode)
	g.Handle("OPTIONS", "", response.CreateOptions("GET", "POST"))

	g.Handle("GET", "/:id", middleware.RequiresPermission(pufferpanel.ScopeNodesView, false), getNode)
	g.Handle("PUT", "/:id", middleware.RequiresPermission(pufferpanel.ScopeNodesEdit, false), updateNode)
	g.Handle("DELETE", "/:id", middleware.RequiresPermission(pufferpanel.ScopeNodesDelete, false), deleteNode)
	g.Handle("OPTIONS", "/:id", response.CreateOptions("PUT", "GET", "DELETE"))

	g.Handle("GET", "/:id/features", middleware.RequiresPermission(pufferpanel.ScopeNodesView, false), getFeatures)
	g.Handle("OPTIONS", "/:id/features", response.CreateOptions("GET"))

	g.Handle("GET", "/:id/deployment", middleware.RequiresPermission(pufferpanel.ScopeNodesDeploy, false), deployNode)
	g.Handle("OPTIONS", "/:id/deployment", response.CreateOptions("GET"))
}

// @Summary Get nodes
// @Description Gets all nodes registered to the panel
// @Success 200 {object} models.NodesView "Nodes"
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Router /api/nodes [get]
// @Security OAuth2Application[nodes.view]
func getAllNodes(c *gin.Context) {
	var err error
	db := middleware.GetDatabase(c)
	ns := &services.Node{DB: db}

	var nodes []*models.Node
	if nodes, err = ns.GetAll(); response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	data := models.FromNodes(nodes)
	c.JSON(http.StatusOK, data)
}

// @Summary Get node
// @Description Gets information about a single node
// @Success 200 {object} models.NodeView "Nodes"
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Param id path string true "Node Id"
// @Router /api/nodes/{id} [get]
// @Security OAuth2Application[nodes.view]
func getNode(c *gin.Context) {
	var err error
	db := middleware.GetDatabase(c)
	ns := &services.Node{DB: db}

	id, ok := validateId(c)
	if !ok {
		return
	}

	node, err := ns.Get(id)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	d := models.FromNode(node)
	c.JSON(http.StatusOK, d)
}

// @Summary Create node
// @Description Creates a node
// @Success 200 {object} models.NodeView "Node created"
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Router /api/nodes [post]
// @Security OAuth2Application[nodes.edit]
func createNode(c *gin.Context) {
	var err error
	db := middleware.GetDatabase(c)
	ns := &services.Node{DB: db}

	model := &models.NodeView{}
	if err = c.BindJSON(model); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	if err = model.Valid(false); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	create := &models.Node{}
	model.CopyToModel(create)

	srt, err := uuid.NewV4()
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	create.Secret = strings.Replace(srt.String(), "-", "", -1)
	if err = ns.Create(create); response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, create)
}

// @Summary Update node
// @Description Updates a node with given information
// @Success 204 {object} nil
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Param id path string true "Node Id"
// @Param node body models.NodeView true "Node information"
// @Router /api/nodes/{id} [put]
// @Security OAuth2Application[nodes.edit]
func updateNode(c *gin.Context) {
	var err error
	db := middleware.GetDatabase(c)
	ns := &services.Node{DB: db}

	viewModel := &models.NodeView{}
	if err = c.BindJSON(viewModel); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	id, ok := validateId(c)
	if !ok {
		return
	}

	if err = viewModel.Valid(true); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	node, err := ns.Get(id)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	viewModel.CopyToModel(node)
	if err = ns.Update(node); response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Deletes a node
// @Description Deletes the node
// @Success 204 {object} nil
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Param id path string true "Node Id"
// @Router /api/nodes/{id} [delete]
// @Security OAuth2Application[nodes.edit]
func deleteNode(c *gin.Context) {
	var err error
	db := middleware.GetDatabase(c)
	ns := &services.Node{DB: db}

	id, ok := validateId(c)
	if !ok {
		return
	}

	node, err := ns.Get(id)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = ns.Delete(node.ID)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Gets the data to deploy a node
// @Description Gets the secret information needed to deploy a node.
// @Success 200 {object} models.Deployment
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Param id path string true "Node Id"
// @Router /api/nodes/{id}/deployment [get]
// @Security OAuth2Application[nodes.deploy]
func deployNode(c *gin.Context) {
	var err error
	db := middleware.GetDatabase(c)
	ns := &services.Node{DB: db}

	id, ok := validateId(c)
	if !ok {
		return
	}

	node, err := ns.Get(id)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	data := &models.Deployment{
		ClientId:     fmt.Sprintf(".node_%d", node.ID),
		ClientSecret: node.Secret,
	}

	c.JSON(http.StatusOK, data)
}

// @Summary Gets the features a node supports
// @Description Gets the environments and if docker is supported on a node
// @Success 200 {object} daemon.Features
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Param id path string true "Node Id"
// @Router /api/nodes/{id}/features [get]
// @Security OAuth2Application[nodes.view]
func getFeatures(c *gin.Context) {
	var err error
	db := middleware.GetDatabase(c)
	ns := &services.Node{DB: db}

	id, ok := validateId(c)
	if !ok {
		return
	}

	node, err := ns.Get(id)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	res, err := ns.CallNode(node, "GET", "/daemon/features", nil, nil)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}
	defer pufferpanel.CloseResponse(res)

	features := &daemon.Features{}
	err = json.NewDecoder(res.Body).Decode(&features)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, features)
}

func validateId(c *gin.Context) (uint, bool) {
	param := c.Param("id")

	id, err := strconv.Atoi(param)

	if response.HandleError(c, err, http.StatusBadRequest) || id < 0 {
		response.HandleError(c, pufferpanel.ErrFieldTooSmall("id", 0), http.StatusBadRequest)
		return 0, false
	}

	return uint(id), true
}
