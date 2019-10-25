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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/v4/response"
	"github.com/pufferpanel/apufferi/v4/scope"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"github.com/pufferpanel/pufferpanel/v2/web/handlers"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func registerNodes(g *gin.RouterGroup) {
	g.Handle("GET", "", handlers.OAuth2Handler(scope.NodesView, false), getAllNodes)
	g.Handle("OPTIONS", "", response.CreateOptions("GET"))
	g.Handle("POST", "", handlers.OAuth2Handler(scope.NodesEdit, false), createNode)

	g.Handle("GET", "/:id", handlers.OAuth2Handler(scope.NodesView, false), getNode)
	g.Handle("PUT", "/:id", handlers.OAuth2Handler(scope.NodesEdit, false), updateNode)
	g.Handle("DELETE", "/:id", handlers.OAuth2Handler(scope.NodesEdit, false), deleteNode)
	g.Handle("OPTIONS", "/:id", response.CreateOptions("PUT", "GET", "POST", "DELETE"))

	g.Handle("GET", "/:id/deployment", handlers.OAuth2Handler(scope.NodesDeploy, false), deployNode)
	g.Handle("OPTIONS", "/:id/deployment", response.CreateOptions("GET"))

	//g.Handle("POST", "/:id/reset", handlers.OAuth2(scope.NodesDeploy, false), response.NotImplemented)
	//g.Handle("OPTIONS", "/:id/reset", response.CreateOptions("POST"))
}

// @Summary Get nodes
// @Description Gets all nodes registered to the panel
// @Accept json
// @Produce json
// @Success 200 {object} models.NodesView "Nodes"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /nodes [get]
func getAllNodes(c *gin.Context) {
	var err error
	db := handlers.GetDatabase(c)
	ns := &services.Node{DB: db}

	var nodes *models.Nodes
	if nodes, err = ns.GetAll(); response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	data := models.FromNodes(nodes)

	c.JSON(http.StatusOK, data)
}

// @Summary Get node
// @Description Gets information about a single node
// @Accept json
// @Produce json
// @Success 200 {object} models.NodeView "Nods"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /nodes/{id} [get]
func getNode(c *gin.Context) {
	var err error
	db := handlers.GetDatabase(c)
	ns := &services.Node{DB: db}

	id, ok := validateId(c)
	if !ok {
		return
	}

	node, err := ns.Get(id)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	data := models.FromNode(node)

	c.JSON(http.StatusOK, data)
}

// @Summary Create node
// @Description Creates a node
// @Accept json
// @Produce json
// @Success 200 {object} models.NodeView "Node created"
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path string true "Node Identifier"
// @Router /nodes [post]
func createNode(c *gin.Context) {
	var err error
	db := handlers.GetDatabase(c)
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
	create.Secret = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	if err = ns.Create(create); response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, create)
}

// @Summary Update node
// @Description Updates a node with given information
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path string true "Node Id"
// @Router /nodes/{id} [put]
func updateNode(c *gin.Context) {
	var err error
	db := handlers.GetDatabase(c)
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
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path string true "Node Id"
// @Router /nodes/{id} [delete]
func deleteNode(c *gin.Context) {
	var err error
	db := handlers.GetDatabase(c)
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

func deployNode(c *gin.Context) {
	var err error
	db := handlers.GetDatabase(c)
	ns := &services.Node{DB: db}

	id, ok := validateId(c)
	if !ok {
		return
	}

	node, err := ns.Get(id)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	services.ValidateTokenLoaded()
	file, err := ioutil.ReadFile(viper.GetString("token.public"))
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	data := &Deployment{
		ClientId:     fmt.Sprintf(".node_%d", node.ID),
		ClientSecret: node.Secret,
		PublicKey:    string(file),
	}

	c.JSON(http.StatusOK, data)
}

func validateId(c *gin.Context) (uint, bool) {
	param := c.Param("id")

	id, err := strconv.Atoi(param)

	if response.HandleError(c, err, http.StatusBadRequest) || id <= 0 {
		response.HandleError(c, pufferpanel.ErrFieldTooSmall("id", 0), http.StatusBadRequest)
		return 0, false
	}

	return uint(id), true
}

type Deployment struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	PublicKey    string `json:"publicKey"`
}
