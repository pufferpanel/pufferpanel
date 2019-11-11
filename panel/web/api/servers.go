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
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/panel/database"
	"github.com/pufferpanel/pufferpanel/v2/panel/models"
	"github.com/pufferpanel/pufferpanel/v2/panel/services"
	"github.com/pufferpanel/pufferpanel/v2/panel/web/handlers"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"github.com/pufferpanel/pufferpanel/v2/scope"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
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
	g.Handle("GET", "/:serverId/user/:email", handlers.OAuth2Handler(scope.ServersEditUsers, true), getServerUsers)
	g.Handle("PUT", "/:serverId/user/:email", handlers.OAuth2Handler(scope.ServersEditUsers, true), handlers.HasTransaction, editServerUser)
	g.Handle("DELETE", "/:serverId/user/:email", handlers.OAuth2Handler(scope.ServersEditUsers, true), handlers.HasTransaction, removeServerUser)
	g.Handle("OPTIONS", "/:serverId", response.CreateOptions("PUT", "GET", "POST", "DELETE"))
}

func searchServers(c *gin.Context) {
	var err error
	db := handlers.GetDatabase(c)
	ss := &services.Server{DB: db}
	ps := &services.Permission{DB: db}

	username := c.DefaultQuery("username", "")
	nodeQuery := c.DefaultQuery("node", "0")
	nameFilter := c.DefaultQuery("name", "*")
	pageSizeQuery := c.DefaultQuery("limit", strconv.Itoa(DefaultPageSize))
	pageQuery := c.DefaultQuery("page", strconv.Itoa(1))

	pageSize, err := strconv.Atoi(pageSizeQuery)
	if response.HandleError(c, err, http.StatusBadRequest) || pageSize <= 0 {
		response.HandleError(c, pufferpanel.ErrFieldTooSmall("pageSize", 0), http.StatusBadRequest)
		return
	}

	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	page, err := strconv.Atoi(pageQuery)
	if response.HandleError(c, err, http.StatusBadRequest) || page <= 0 {
		response.HandleError(c, pufferpanel.ErrFieldTooSmall("page", 0), http.StatusBadRequest)
		return
	}

	node, err := strconv.Atoi(nodeQuery)
	if response.HandleError(c, err, http.StatusBadRequest) || page <= 0 {
		response.HandleError(c, pufferpanel.ErrFieldTooSmall("nodeId", 0), http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(*models.User)

	perms, err := ps.GetForUser(user.ID)
	if response.HandleError(c, err, http.StatusInternalServerError) {
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
		c.JSON(http.StatusOK, &ServerSearchResponse{
			Servers: []*models.ServerView{},
			Metadata: &response.Metadata{Paging: &response.Paging{
				Page:    1,
				Size:    0,
				MaxSize: MaxPageSize,
				Total:   0,
			}},
		})
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
	if results, total, err = ss.Search(searchCriteria); response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, &ServerSearchResponse{
		Servers: models.RemoveServerPrivateInfoFromAll(models.FromServers(results)),
		Metadata: &response.Metadata{Paging: &response.Paging{
			Page:    uint(page),
			Size:    uint(pageSize),
			MaxSize: MaxPageSize,
			Total:   total,
		}},
	})
}

func getServer(c *gin.Context) {
	t, exist := c.Get("server")

	if !exist {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	server, ok := t.(*models.Server)
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
	}

	_, includePerms := c.GetQuery("perms")
	var perms *models.PermissionView
	if includePerms {
		db, err := database.GetConnection()
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}

		u := c.MustGet("user").(*models.User)

		ps := &services.Permission{DB: db}

		p, err := ps.GetForUserAndServer(u.ID, &server.Identifier)
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}
		perms = models.FromPermission(p)
	}

	data := &GetServerResponse{
		Server: models.RemoveServerPrivateInfo(models.FromServer(server)),
		Perms:  perms,
	}

	c.JSON(http.StatusOK, data)
}

func createServer(c *gin.Context) {
	var err error
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
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	node, err := ns.Get(postBody.NodeId)

	if gorm.IsRecordNotFoundError(err) {
		response.HandleError(c, pufferpanel.ErrNodeInvalid, http.StatusBadRequest)
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	port, err := getFromDataOrDefault(postBody.Variables, "port", uint16(0))
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	ip, err := getFromDataOrDefault(postBody.Variables, "ip", "0.0.0.0")
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	if postBody.Name == "" {
		postBody.Name = postBody.Identifier
	}

	server := &models.Server{
		Name:       postBody.Name,
		Identifier: postBody.Identifier,
		NodeID:     node.ID,
		IP:         ip.(string),
		Port:       port.(uint16),
		Type:       postBody.Type,
	}

	users := make([]*models.User, len(postBody.Users))

	for k, v := range postBody.Users {
		user, err := us.Get(v)
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}

		users[k] = user
	}

	err = ss.Create(server)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	for _, v := range users {
		perm, err := ps.GetForUserAndServer(v.ID, &server.Identifier)
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}

		perm.SetDefaults()

		err = ps.UpdatePermissions(perm)
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}
	}

	data, _ := json.Marshal(postBody.Server)
	reader := ioutil.NopCloser(bytes.NewReader(data))

	//we need to get your new token
	token, err := ps.GenerateOAuthForUser(c.MustGet("user").(*models.User).ID, nil)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)

	nodeResponse, err := ns.CallNode(node, "PUT", "/server/"+server.Identifier, reader, headers)
	if nodeResponse != nil {
		defer pufferpanel.Close(nodeResponse.Body)
	}

	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	if nodeResponse.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(nodeResponse.Body)
		logging.Error().Printf("Unexpected response from daemon: %+v\n%s", nodeResponse.StatusCode, buf.String())
		response.HandleError(c, pufferpanel.ErrUnknownError, http.StatusInternalServerError)
		return
	}

	if response.HandleError(c, db.Commit().Error, http.StatusInternalServerError) {
		return
	}

	es := services.GetEmailService()
	for _, user := range users {
		err = es.SendEmail(user.Email, "addedToServer", map[string]interface{}{
			"Server":        server,
			"RegisterToken": "",
		}, true)
		if err != nil {
			//since we don't want to tell the user it failed, we'll log and move on
			logging.Error().Printf("Error sending email: %s", err)
		}
	}

	c.JSON(http.StatusOK, &CreateServerResponse{Id: serverId})
}

func deleteServer(c *gin.Context) {
	var err error

	db := handlers.GetDatabase(c)
	ss := &services.Server{DB: db}

	t, exist := c.Get("server")

	if !exist {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	server, ok := t.(*models.Server)
	if !ok {
		response.HandleError(c, pufferpanel.ErrUnknownError, http.StatusInternalServerError)
		return
	}

	err = ss.Delete(server.Identifier)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	if response.HandleError(c, db.Commit().Error, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

func getServerUsers(c *gin.Context) {
	var err error
	db := handlers.GetDatabase(c)
	ps := &services.Permission{DB: db}

	t, exist := c.Get("server")

	if !exist {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	server, ok := t.(*models.Server)
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	perms, err := ps.GetForServer(server.Identifier)
	if response.HandleError(c, err, http.StatusInternalServerError) {
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

	c.JSON(http.StatusOK, users)
}

func editServerUser(c *gin.Context) {
	var err error
	db := handlers.GetDatabase(c)
	us := &services.User{DB: db}
	ps := &services.Permission{DB: db}

	email := c.Param("email")
	if email == "" {
		return
	}

	perms := &models.PermissionView{}
	err = c.BindJSON(perms)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}
	perms.Email = email

	t, exist := c.Get("server")
	if !exist {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	server, ok := t.(*models.Server)
	if !ok {
		response.HandleError(c, pufferpanel.ErrServerNotFound, http.StatusInternalServerError)
		return
	}

	var registerToken string
	user, err := us.GetByEmail(email)
	if err != nil && !gorm.IsRecordNotFoundError(err) && response.HandleError(c, err, http.StatusInternalServerError) {
		return
	} else if gorm.IsRecordNotFoundError(err) {
		//we need to create the user here, since it's a new email we've not seen

		user = &models.User{
			Username: uuid.NewV4().String(),
			Email:    email,
		}
		registerToken = uuid.NewV4().String()
		err = user.SetPassword(registerToken)
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}

		err = us.Create(user)
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}
	}

	existing, err := ps.GetForUserAndServer(user.ID, &server.Identifier)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	var firstTimeAccess = false
	if existing.ID == 0 {
		firstTimeAccess = true
	}
	perms.CopyTo(existing, false)
	err = ps.UpdatePermissions(existing)

	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	if response.HandleError(c, db.Commit().Error, http.StatusInternalServerError) {
		return
	}

	//now we can send emails to the people
	if firstTimeAccess {
		es := services.GetEmailService()
		err = es.SendEmail(user.Email, "addedToServer", map[string]interface{}{
			"Server":        server,
			"RegisterToken": registerToken,
		}, true)
		if err != nil {
			//since we don't want to tell the user it failed, we'll log and move on
			logging.Error().Printf("Error sending email: %s", err)
		}
	}

	c.Status(http.StatusNoContent)
}

func removeServerUser(c *gin.Context) {
	var err error
	db := handlers.GetDatabase(c)
	us := &services.User{DB: db}
	ps := &services.Permission{DB: db}

	email := c.Param("email")
	if email == "" {
		return
	}

	t, exist := c.Get("server")

	if !exist {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	server, ok := t.(*models.Server)
	if !ok {
		response.HandleError(c, pufferpanel.ErrServerNotFound, http.StatusInternalServerError)
		return
	}

	user, err := us.GetByEmail(email)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	perms, err := ps.GetForUserAndServer(user.ID, &server.Identifier)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = ps.Remove(perms)

	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	if response.HandleError(c, db.Commit().Error, http.StatusInternalServerError) {
		return
	}
}

type serverCreation struct {
	pufferpanel.Server

	NodeId uint     `json:"node,string"`
	Users  []string `json:"users"`
	Name   string   `json:"name"`
}

func getFromData(variables map[string]pufferpanel.Variable, key string) interface{} {
	for k, v := range variables {
		if k == key {
			return v.Value
		}
	}
	return nil
}

func getFromDataOrDefault(variables map[string]pufferpanel.Variable, key string, val interface{}) (interface{}, error) {
	res := getFromData(variables, key)

	if res != nil {
		return pufferpanel.Convert(res, val)
	}

	return val, nil
}

type GetServerResponse struct {
	Server *models.ServerView     `json:"server"`
	Perms  *models.PermissionView `json:"permissions"`
}

type CreateServerResponse struct {
	Id string `json:"id"`
}

type ServerSearchResponse struct {
	Servers []*models.ServerView `json:"servers"`
	*response.Metadata
}
