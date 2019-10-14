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
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/apufferi/v3"
	"github.com/pufferpanel/apufferi/v3/logging"
	"github.com/pufferpanel/apufferi/v3/response"
	"github.com/pufferpanel/apufferi/v3/scope"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"github.com/pufferpanel/pufferpanel/v2/web/handlers"
	"github.com/satori/go.uuid"
	"io"
	"net/http"
	"reflect"
	"strconv"
)

func registerServers(g *gin.RouterGroup) {
	g.Handle("GET", "", handlers.OAuth2Handler(scope.ServersView, false), searchServers)
	g.Handle("OPTIONS", "", response.CreateOptions("GET"))

	g.Handle("POST", "", handlers.OAuth2Handler(scope.ServersCreate, false), handlers.HasTransaction, createServer)
	g.Handle("GET", "/:serverId", handlers.OAuth2Handler(scope.ServersView, true), getServer)
	g.Handle("PUT", "/:serverId", handlers.OAuth2Handler(scope.ServersCreate, false), handlers.HasTransaction, createServer)
	g.Handle("POST", "/:serverId", handlers.OAuth2Handler(scope.ServersEdit, true), handlers.HasTransaction, createServer)
	g.Handle("DELETE", "/:serverId", handlers.OAuth2Handler(scope.ServersDelete, true), handlers.HasTransaction, deleteServer)
	g.Handle("GET", "/:serverId/user", handlers.OAuth2Handler(scope.ServersEditUsers, true), getServerUsers)
	g.Handle("GET", "/:serverId/user/:username", handlers.OAuth2Handler(scope.ServersEditUsers, true), getServerUsers)
	g.Handle("PUT", "/:serverId/user/:username", handlers.OAuth2Handler(scope.ServersEditUsers, true), handlers.HasTransaction, editServerUser)
	g.Handle("DELETE", "/:serverId/user/:username", handlers.OAuth2Handler(scope.ServersEditUsers, true), handlers.HasTransaction, removeServerUser)
	g.Handle("OPTIONS", "/:serverId", response.CreateOptions("PUT", "GET", "POST", "DELETE"))
}

func searchServers(c *gin.Context) {
	var err error
	res := response.From(c)
	db := handlers.GetDatabase(c)
	ss := &services.Server{DB: db}
	ps := &services.Permission{DB: db}

	username := c.DefaultQuery("username", "")
	nodeQuery := c.DefaultQuery("node", "0")
	nameFilter := c.DefaultQuery("name", "*")
	pageSizeQuery := c.DefaultQuery("limit", strconv.Itoa(DefaultPageSize))
	pageQuery := c.DefaultQuery("page", strconv.Itoa(1))

	pageSize, err := strconv.Atoi(pageSizeQuery)
	if err != nil || pageSize <= 0 {
		res.Fail().Status(http.StatusBadRequest).Message("page size must be a positive number")
		return
	}

	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	page, err := strconv.Atoi(pageQuery)
	if err != nil || page <= 0 {
		res.Fail().Status(http.StatusBadRequest).Message("page must be a positive number")
		return
	}

	node, err := strconv.Atoi(nodeQuery)
	if err != nil || page <= 0 {
		res.Fail().Status(http.StatusBadRequest).Message("node id is invalid")
		return
	}

	user := c.MustGet("user").(*models.User)

	perms, err := ps.GetForUser(user.ID)
	if response.HandleError(res, err) {
		return
	}

	isAdmin := false
	for _, p := range perms {
		if p.Admin {
			isAdmin = true
			break
		}
	}

	if !isAdmin && username != "" && user.Username != username {
		res.PageInfo(uint(page), uint(pageSize), MaxPageSize, 0).Data(make([]models.ServerView, 0))
		return
	} else if !isAdmin {
		username = user.Username
	}

	var results *models.Servers
	var total uint
	searchCriteria := services.ServerSearch{
		Username: username,
		NodeId:   uint(node),
		Name:     nameFilter,
		PageSize: uint(pageSize),
		Page:     uint(page),
	}
	if results, total, err = ss.Search(searchCriteria); response.HandleError(res, err) {
		return
	}

	res.PageInfo(uint(page), uint(pageSize), MaxPageSize, total).Data(models.RemoveServerPrivateInfoFromAll(models.FromServers(results)))
}

func getServer(c *gin.Context) {
	res := response.From(c)

	t, exist := c.Get("server")

	if !exist {
		response.HandleError(res, pufferpanel.ErrServerNotFound)
		return
	}

	server, ok := t.(*models.Server)
	if !ok {
		response.HandleError(res, pufferpanel.ErrServerNotFound)
	}

	_, includePerms := c.GetQuery("perms")
	var perms *models.PermissionView
	if includePerms {
		db, err := database.GetConnection()
		if response.HandleError(res, err) {
			return
		}

		u := c.MustGet("user").(*models.User)

		ps := &services.Permission{DB: db}

		p, err := ps.GetForUserAndServer(u.ID, &server.Identifier)
		if response.HandleError(res, err) {
			return
		}
		perms = models.FromPermission(p)
	}

	data := &GetServerResponse{
		Server: models.RemoveServerPrivateInfo(models.FromServer(server)),
		Perms:  perms,
	}

	res.Data(data)
}

func createServer(c *gin.Context) {
	var err error
	res := response.From(c)
	db := handlers.GetDatabase(c)
	ss := &services.Server{DB: db}
	ns := &services.Node{DB: db}
	us := &services.User{DB: db}
	ps := &services.Permission{DB: db}

	serverId := c.Param("id")
	if serverId == "" {
		serverId = uuid.NewV4().String()[:8]
	}

	postBody := &serverCreation{}
	err = c.Bind(postBody)
	postBody.Identifier = serverId
	if err != nil {
		res.Status(http.StatusBadRequest).Error(err).Fail()
		return
	}

	node, exists, err := ns.Get(postBody.NodeId)

	if response.HandleError(res, err) {
		return
	}

	if !exists {
		res.Status(http.StatusBadRequest).Message("no node with given id").Fail()
	}

	server := &models.Server{
		Name:       getFromDataOrDefault(postBody.Variables, "name", postBody.Identifier).(string),
		Identifier: postBody.Identifier,
		NodeID:     node.ID,
		IP:         getFromDataOrDefault(postBody.Variables, "ip", "0.0.0.0").(string),
		Port:       getFromDataOrDefault(postBody.Variables, "port", uint(0)).(uint),
		Type:       postBody.Type,
	}

	users := make([]*models.User, len(postBody.Users))

	for k, v := range postBody.Users {
		user, err := us.Get(v)
		if response.HandleError(res, err) {
			return
		}
		if !exists {
			response.HandleError(res, pufferpanel.ErrUserNotFound.Metadata(map[string]interface{}{"username": v}))
			return
		}

		users[k] = user
	}

	err = ss.Create(server)
	if err != nil {
		res.Status(http.StatusInternalServerError).Error(err).Fail()
		return
	}

	for _, v := range users {
		perm, err := ps.GetForUserAndServer(v.ID, &server.Identifier)
		if response.HandleError(res, err) {
			return
		}

		perm.SetDefaults()

		err = ps.UpdatePermissions(perm)
		if response.HandleError(res, err) {
			return
		}
	}

	data, _ := json.Marshal(postBody.Server)
	reader := newFakeReader(data)

	headers := http.Header{}
	headers.Set("Authorization", c.GetHeader("Authorization"))

	nodeResponse, err := ns.CallNode(node, "PUT", "/server/"+server.Identifier, reader, headers)

	if response.HandleError(res, err) {
		return
	}

	if nodeResponse.StatusCode != http.StatusOK {
		logging.Build(logging.ERROR).WithMessage("Unexpected response from daemon: %+v").WithArgs(nodeResponse.StatusCode).Log()
		response.HandleError(res, pufferpanel.ErrUnknownError)
		return
	}

	apiResponse := &response.Response{}
	err = json.NewDecoder(nodeResponse.Body).Decode(apiResponse)

	if response.HandleError(res, err) {
		return
	}

	if !apiResponse.Success {
		logging.Build(logging.ERROR).WithMessage("Unexpected response from daemon: %+v").WithArgs(apiResponse).Log()
		response.HandleError(res, pufferpanel.ErrUnknownError)
		return
	}

	if response.HandleError(res, db.Commit().Error) {
		return
	}

	res.Data(server.Identifier)
}

func deleteServer(c *gin.Context) {
	var err error
	res := response.From(c)

	db, err := database.GetConnection()
	if response.HandleError(res, err) {
		return
	}
	db.Begin()

	ss := &services.Server{DB: db}

	t, exist := c.Get("server")

	if !exist {
		response.HandleError(res, pufferpanel.ErrServerNotFound)
		return
	}

	server, ok := t.(*models.Server)
	if !ok {
		response.HandleError(res, pufferpanel.ErrServerNotFound)
		return
	}

	err = ss.Delete(server.Identifier)
	if response.HandleError(res, err) {
		return
	}

	if response.HandleError(res, db.Commit().Error) {
		return
	}

	v := models.FromServer(server)
	res.Status(http.StatusOK).Data(v)
}

func getServerUsers(c *gin.Context) {
	var err error
	res := response.From(c)
	db := handlers.GetDatabase(c)
	ps := &services.Permission{DB: db}

	t, exist := c.Get("server")

	if !exist {
		response.HandleError(res, pufferpanel.ErrServerNotFound)
		return
	}

	server, ok := t.(*models.Server)
	if !ok {
		response.HandleError(res, pufferpanel.ErrServerNotFound)
		return
	}

	perms, err := ps.GetForServer(server.Identifier)
	if response.HandleError(res, err) {
		return
	}

	caller, _ := c.MustGet("user").(*models.User)

	for i := 0; i < len(perms); i++ {
		if *perms[i].UserId == caller.ID {
			perms = append(perms[:i], perms[i+1:]...)
			i--
		}
	}

	users := make([]*models.PermissionView, len(perms))
	for k, v := range perms {
		users[k] = models.FromPermission(v)
	}

	res.Data(users)
}

func editServerUser(c *gin.Context) {
	var err error
	res := response.From(c)
	db := handlers.GetDatabase(c)
	us := &services.User{DB: db}
	ps := &services.Permission{DB: db}

	username := c.Param("username")
	if username == "" {
		return
	}

	perms := &models.PermissionView{}
	err = c.BindJSON(perms)
	if response.HandleError(res, err) {
		return
	}
	perms.Username = username

	t, exist := c.Get("server")

	if !exist {
		response.HandleError(res, pufferpanel.ErrServerNotFound)
		return
	}

	server, ok := t.(*models.Server)
	if !ok {
		response.HandleError(res, pufferpanel.ErrServerNotFound)
		return
	}

	user, err := us.Get(username)
	if response.HandleError(res, err) {
		return
	}

	existing, err := ps.GetForUserAndServer(user.ID, &server.Identifier)
	if response.HandleError(res, err) {
		return
	}
	perms.CopyTo(existing)
	err = ps.UpdatePermissions(existing)

	if response.HandleError(res, err) {
		return
	}

	if response.HandleError(res, db.Commit().Error) {
		return
	}
}

func removeServerUser(c *gin.Context) {
	var err error
	res := response.From(c)
	db := handlers.GetDatabase(c)
	us := &services.User{DB: db}
	ps := &services.Permission{DB: db}

	username := c.Param("username")
	if username == "" {
		return
	}

	t, exist := c.Get("server")

	if !exist {
		response.HandleError(res, pufferpanel.ErrServerNotFound)
		return
	}

	server, ok := t.(*models.Server)
	if !ok {
		response.HandleError(res, pufferpanel.ErrServerNotFound)
		return
	}

	user, err := us.Get(username)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		res.Fail().Status(http.StatusNotFound).Message("no user with username")
		return
	} else if response.HandleError(res, err) {
		return
	}

	perms, err := ps.GetForUserAndServer(user.ID, &server.Identifier)
	if response.HandleError(res, err) {
		return
	}

	err = ps.Remove(perms)

	if response.HandleError(res, err) {
		return
	}

	if response.HandleError(res, db.Commit().Error) {
		return
	}
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

type serverCreation struct {
	apufferi.Server

	NodeId uint     `json:"node"`
	Users  []string `json:"users"`
}

func getFromData(variables map[string]apufferi.Variable, key string) interface{} {
	for k, v := range variables {
		if k == key {
			return v.Value
		}
	}
	return nil
}

//this will enforce whatever the type val is defined as will be what is returned
func getFromDataOrDefault(variables map[string]apufferi.Variable, key string, val interface{}) interface{} {
	res := getFromData(variables, key)

	if res != nil {
		if reflect.TypeOf(val).AssignableTo(reflect.TypeOf(res)) {
			return res
		}
	}

	return val
}

type GetServerResponse struct {
	Server *models.ServerView     `json:"server"`
	Perms  *models.PermissionView `json:"permissions"`
}
