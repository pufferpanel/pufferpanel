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
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi"
	builder "github.com/pufferpanel/apufferi/response"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/errors"
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/pufferpanel/pufferpanel/models/view"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/shared"
	"github.com/pufferpanel/pufferpanel/web/handlers"
	"github.com/satori/go.uuid"
	"io"
	"net/http"
	"strconv"
)

func registerServers(g *gin.RouterGroup) {
	g.Handle("GET", "", handlers.OAuth2WithLimit("servers.view", false), searchServers)
	g.Handle("OPTIONS", "", shared.CreateOptions("GET"))

	g.Handle("POST", "", handlers.OAuth2("servers.create", false), createServer)
	g.Handle("GET", "/:serverId", handlers.OAuth2("servers.view", true), getServer)
	g.Handle("PUT", "/:serverId", handlers.OAuth2("servers.edit", false), createServer)
	g.Handle("DELETE", "/:serverId", handlers.OAuth2("servers.edit", false), deleteServer)
	g.Handle("GET", "/:serverId/users", handlers.OAuth2("servers.edit", true), getServerUsers)
	g.Handle("POST", "/:serverId/users", handlers.OAuth2("servers.edit", true), editServerUsers)
	g.Handle("OPTIONS", "/:serverId", shared.CreateOptions("PUT", "GET", "POST", "DELETE"))
}

func searchServers(c *gin.Context) {
	var err error
	response := builder.From(c)

	username := c.DefaultQuery("username", "")
	nodeQuery := c.DefaultQuery("node", "0")
	nameFilter := c.DefaultQuery("name", "*")
	pageSizeQuery := c.DefaultQuery("limit", strconv.Itoa(DefaultPageSize))
	pageQuery := c.DefaultQuery("page", strconv.Itoa(1))

	pageSize, err := strconv.Atoi(pageSizeQuery)
	if err != nil || pageSize <= 0 {
		response.Fail().Status(http.StatusBadRequest).Message("page size must be a positive number")
		return
	}

	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	page, err := strconv.Atoi(pageQuery)
	if err != nil || page <= 0 {
		response.Fail().Status(http.StatusBadRequest).Message("page must be a positive number")
		return
	}

	node, err := strconv.Atoi(nodeQuery)
	if err != nil || page <= 0 {
		response.Fail().Status(http.StatusBadRequest).Message("node id is invalid")
		return
	}

	db, err := database.GetConnection()
	if shared.HandleError(response, err) {
		return
	}

	ss := &services.Server{DB: db}
	os := services.GetOAuth(db)

	//see if user has access to view all others, otherwise we can't permit search without their username
	ci, allowed, _ := os.HasRights(c.GetString("accessToken"), nil, "servers.view");
	if !allowed {
		response.PageInfo(uint(page), uint(pageSize), MaxPageSize, 0).Data(make([]view.ServerViewModel, 0))
		return
	}

	username = ci.User.Username

	var results *models.Servers
	var total uint
	searchCriteria := services.ServerSearch{
		Username: username,
		NodeId:   uint(node),
		Name:     nameFilter,
		PageSize: uint(pageSize),
		Page:     uint(page),
	}
	if results, total, err = ss.Search(searchCriteria); shared.HandleError(response, err) {
		return
	}

	response.PageInfo(uint(page), uint(pageSize), MaxPageSize, total).Data(view.RemoveServerPrivateInfoFromAll(view.FromServers(results)))
}

func getServer(c *gin.Context) {
	response := builder.From(c)

	t, exist := c.Get("server")

	if !exist {
		shared.HandleError(response, errors.ErrServerNotFound)
		return
	}

	server, ok := t.(*models.Server)
	if !ok {
		shared.HandleError(response, errors.ErrServerNotFound)
	}

	response.Data(view.RemoveServerPrivateInfo(view.FromServer(server)))
}

func createServer(c *gin.Context) {
	var err error
	response := builder.From(c)

	serverId := c.Param("id")
	if serverId == "" {
		serverId = uuid.NewV4().String()[:8]
	}

	postBody := view.ServerViewModel{}
	err = c.Bind(&postBody)
	postBody.Identifier = serverId
	if err != nil {
		response.Status(http.StatusBadRequest).Error(err).Fail()
		return
	}

	db, err := database.GetConnection()
	if shared.HandleError(response, err) {
		return
	}

	//time for a transaction!
	trans := db.Begin()
	success := false
	defer func() {
		if !success {
			trans.Rollback()
		}
	}()

	ss := &services.Server{DB: trans}
	ns := &services.Node{DB: trans}

	node, exists, err := ns.Get(postBody.NodeId)

	if shared.HandleError(response, err) {
		return
	}

	if !exists {
		response.Status(http.StatusBadRequest).Message("no node with given id").Fail()
	}

	server := &models.Server{}
	postBody.CopyToModel(server)

	server.NodeID = node.ID

	err = ss.Create(server)
	if err != nil {
		response.Status(http.StatusInternalServerError).Error(err).Fail()
		return
	}

	apiServer := apufferi.Server{
		Variables:      nil,
		Display:        "",
		Environment:    apufferi.TypeWithMetadata{},
		Installation:   nil,
		Uninstallation: nil,
		Type:           "",
		Identifier:     "",
		Execution:      apufferi.Execution{},
	}

	data, _ := json.Marshal(apiServer)
	reader := newFakeReader(data)
	nodeResponse, err := ns.CallNode(node, "PUT", "/server/"+server.Identifier, reader, nil)

	if shared.HandleError(response, err) {
		return
	}

	if nodeResponse.StatusCode != http.StatusOK {
		return
	}

	postBody.Data = nil
	response.Data(node)

	trans.Commit()
	success = true
}

func deleteServer(c *gin.Context) {
	var err error
	response := builder.From(c)

	db, err := database.GetConnection()
	if shared.HandleError(response, err) {
		return
	}

	ss := &services.Server{DB: db}

	t, exist := c.Get("server")

	if !exist {
		shared.HandleError(response, errors.ErrServerNotFound)
		return
	}

	server, ok := t.(*models.Server)
	if !ok {
		shared.HandleError(response, errors.ErrServerNotFound)
	}

	err = ss.Delete(server.ID)
	if shared.HandleError(response, err) {
		return
	} else {
		v := view.FromServer(server)
		response.Status(http.StatusOK).Data(v)
	}
}

func getServerUsers(c *gin.Context) {

}

func editServerUsers(c *gin.Context) {

}

//This class exists
type fakeReader struct {
	reader io.Reader
}

func newFakeReader(data []byte) *fakeReader {
	return &fakeReader{reader: bytes.NewReader(data)}
}

func (fr *fakeReader) Read(p []byte) (int, error) {
	return fr.reader.Read(p)
}

func (fr *fakeReader) Close() error {
	return nil
}
