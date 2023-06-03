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
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/middleware"
	"github.com/pufferpanel/pufferpanel/v2/middleware/handlers"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func registerServers(g *gin.RouterGroup) {
	g.Handle("GET", "", handlers.OAuth2Handler(pufferpanel.ScopeServersView, false), searchServers)
	g.Handle("OPTIONS", "", response.CreateOptions("GET"))

	g.Handle("POST", "", handlers.OAuth2Handler(pufferpanel.ScopeServersCreate, false), middleware.HasTransaction, createServer)
	g.Handle("GET", "/:serverId", handlers.OAuth2Handler(pufferpanel.ScopeServersView, true), getServer)
	g.Handle("PUT", "/:serverId", handlers.OAuth2Handler(pufferpanel.ScopeServersCreate, false), middleware.HasTransaction, createServer)
	g.Handle("POST", "/:serverId", handlers.OAuth2Handler(pufferpanel.ScopeServersEdit, true), middleware.HasTransaction, createServer)
	g.Handle("DELETE", "/:serverId", handlers.OAuth2Handler(pufferpanel.ScopeServersDelete, true), middleware.HasTransaction, deleteServer)
	g.Handle("PUT", "/:serverId/name/:name", handlers.OAuth2Handler(pufferpanel.ScopeServersEdit, true), middleware.HasTransaction, renameServer)
	g.Handle("OPTIONS", "/:serverId", response.CreateOptions("PUT", "GET", "POST", "DELETE"))

	g.Handle("GET", "/:serverId/user", handlers.OAuth2Handler(pufferpanel.ScopeServersEditUsers, true), getServerUsers)
	g.Handle("OPTIONS", "/:serverId/user", response.CreateOptions("GET"))

	g.Handle("GET", "/:serverId/user/:email", handlers.OAuth2Handler(pufferpanel.ScopeServersEditUsers, true), getServerUsers)
	g.Handle("PUT", "/:serverId/user/:email", handlers.OAuth2Handler(pufferpanel.ScopeServersEditUsers, true), middleware.HasTransaction, editServerUser)
	g.Handle("DELETE", "/:serverId/user/:email", handlers.OAuth2Handler(pufferpanel.ScopeServersEditUsers, true), middleware.HasTransaction, removeServerUser)
	g.Handle("OPTIONS", "/:serverId/user/:email", response.CreateOptions("GET", "PUT", "DELETE"))

	g.Handle("GET", "/:serverId/oauth2", handlers.OAuth2Handler(pufferpanel.ScopeServersView, true), getOAuth2Clients)
	g.Handle("POST", "/:serverId/oauth2", handlers.OAuth2Handler(pufferpanel.ScopeServersView, true), createOAuth2Client)
	g.Handle("OPTIONS", "/:serverId/oauth2", response.CreateOptions("GET", "POST"))

	g.Handle("DELETE", "/:serverId/oauth2/:clientId", handlers.OAuth2Handler(pufferpanel.ScopeServersView, true), deleteOAuth2Client)
	g.Handle("OPTIONS", "/:serverId/oauth2/:clientId", response.CreateOptions("DELETE"))
}

// @Summary Get servers
// @Description Gets servers, and allowing for filtering of servers. * is a wildcard that can be used for text inputs
// @Accept json
// @Produce json
// @Success 200 {object} models.ServerSearchResponse
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param username query string false "Username to filter on, default is current user if NOT admin"
// @Param node query uint false "Node ID to filter on"
// @Param name query string false "Name of server to filter on"
// @Param limit query uint false "Max number of results to return"
// @Param page query uint false "What page to get back for many results"
// @Router /api/servers [get]
func searchServers(c *gin.Context) {
	var err error
	db := middleware.GetDatabase(c)
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
		c.JSON(http.StatusOK, &models.ServerSearchResponse{
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

	searchCriteria := services.ServerSearch{
		Username: username,
		NodeId:   uint(node),
		Name:     nameFilter,
		PageSize: uint(pageSize),
		Page:     uint(page),
	}

	results, total, err := ss.Search(searchCriteria)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	data := models.FromServers(results)
	c.JSON(http.StatusOK, &models.ServerSearchResponse{
		Servers: models.RemoveServerPrivateInfoFromAll(data),
		Metadata: &response.Metadata{Paging: &response.Paging{
			Page:    uint(page),
			Size:    uint(pageSize),
			MaxSize: MaxPageSize,
			Total:   total,
		}},
	})
}

// @Summary Get a server
// @Description Gets a particular server
// @Accept json
// @Produce json
// @Success 200 {object} models.GetServerResponse
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path string true "Server ID"
// @Router /api/servers/{id} [get]
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

	d := &models.GetServerResponse{
		Server: models.RemoveServerPrivateInfo(models.FromServer(server)),
		Perms:  perms,
	}

	if d.Server.Node.PrivateHost == "127.0.0.1" && d.Server.Node.PublicHost == "127.0.0.1" {
		d.Server.Node.PublicHost = strings.SplitN(c.Request.Host, ":", 2)[0]
	}

	c.JSON(http.StatusOK, d)
}

// @Summary Makes a server
// @Description Creates a server
// @Accept json
// @Produce json
// @Success 200 {object} models.CreateServerResponse
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path string false "Server ID"
// @Param server body models.ServerCreation true "Creation information"
// @Router /api/servers [post]
// @Router /api/servers/{id} [put]
func createServer(c *gin.Context) {
	var err error
	db := middleware.GetDatabase(c)
	ss := &services.Server{DB: db}
	ns := &services.Node{DB: db}
	us := &services.User{DB: db}
	ps := &services.Permission{DB: db}

	serverId := c.Param("serverId")

	if serverId == "" {
		serverId = uuid.NewV4().String()[:8]
	}

	postBody := &models.ServerCreation{}
	err = c.Bind(postBody)
	postBody.Identifier = serverId
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	node, err := ns.Get(postBody.NodeId)

	if gorm.ErrRecordNotFound == err {
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
		Type:       postBody.Type.Type,
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

	nodeResponse, err := ns.CallNode(node, "PUT", "/daemon/server/"+server.Identifier, reader, headers)
	if nodeResponse != nil {
		defer pufferpanel.Close(nodeResponse.Body)
	}

	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	if nodeResponse.StatusCode != http.StatusOK {
		resData, err := ioutil.ReadAll(nodeResponse.Body)
		if err != nil {
			logging.Error.Printf("Failed to parse response from daemon\n%s", err.Error())
		}
		logging.Error.Printf("Unexpected response from daemon: %+v\n%s", nodeResponse.StatusCode, string(resData))
		//assume daemon gives us a valid response, directly forward to client
		c.Header("Content-Type", "application/json")
		c.Status(nodeResponse.StatusCode)
		_, _ = c.Writer.Write(resData)
		c.Abort()
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
			logging.Error.Printf("Error sending email: %s", err)
		}
	}

	c.JSON(http.StatusOK, &models.CreateServerResponse{Id: serverId})
}

// @Summary Deletes a server
// @Description Deletes a server from the panel
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path string true "Server ID"
// @Router /api/servers/{id} [delete]
func deleteServer(c *gin.Context) {
	var err error

	db := middleware.GetDatabase(c)
	ss := &services.Server{DB: db}
	ns := &services.Node{DB: db}

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

	t, exist = c.Get("user")
	if !exist {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	user, ok := t.(*models.User)
	if !ok {
		response.HandleError(c, pufferpanel.ErrUnknownError, http.StatusInternalServerError)
		return
	}

	node, err := ns.Get(server.NodeID)
	if err != nil {
		response.HandleError(c, err, http.StatusInternalServerError)
		return
	}

	//we need to know what users are impacted by a server being deleted
	ps := services.Permission{DB: db}
	users := make([]models.User, 0)
	perms, err := ps.GetForServer(server.Identifier)
	for _, p := range perms {
		exists := false
		for _, u := range users {
			if u.ID == p.User.ID {
				exists = true
				break
			}
		}
		if exists {
			continue
		}
		users = append(users, p.User)
	}

	err = ss.Delete(server.Identifier)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	newHeader, err := ps.GenerateOAuthForUser(user.ID, &server.Identifier)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+newHeader)

	nodeRes, err := ns.CallNode(node, "DELETE", "/daemon/server/"+server.Identifier, nil, headers)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		//node didn't permit it, REVERT!
		db.Rollback()
		return
	}

	if nodeRes.StatusCode != http.StatusNoContent {
		response.HandleError(c, errors.New("invalid status code response: "+nodeRes.Status), http.StatusInternalServerError)
		return
	}

	if response.HandleError(c, db.Commit().Error, http.StatusInternalServerError) {
		return
	}

	es := services.GetEmailService()
	for _, u := range users {
		err = es.SendEmail(u.Email, "deletedServer", map[string]interface{}{
			"Server": server,
		}, true)
		if err != nil {
			//since we don't want to tell the user it failed, we'll log and move on
			logging.Error.Printf("Error sending email: %s\n", err)
		}
	}

	c.Status(http.StatusNoContent)
}

// @Summary Gets all users for a server
// @Accept json
// @Produce json
// @Success 200 {object} []models.PermissionView
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path string true "Server ID"
// @Router /api/servers/{id}/user [get]
func getServerUsers(c *gin.Context) {
	var err error
	db := middleware.GetDatabase(c)
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

// @Summary Edits access to a server
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path string true "Server ID"
// @Param email path string true "Email of user"
// @Param body body models.PermissionView true "New permissions to apply"
// @Router /api/servers/{id}/users/{email} [put]
func editServerUser(c *gin.Context) {
	var err error
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}
	ps := &services.Permission{DB: db}

	email := c.Param("email")
	username := c.Param("username")
	if email == "" && username == "" {
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
	var user *models.User
	if email != "" {
		user, err = us.GetByEmail(email)
	} else {
		user, err = us.Get(username)
	}

	if err != nil && gorm.ErrRecordNotFound != err && response.HandleError(c, err, http.StatusInternalServerError) {
		return
	} else if gorm.ErrRecordNotFound == err {
		if email == "" {
			response.HandleError(c, err, http.StatusBadRequest)
			return
		}
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
	existing.ViewServer = true
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
			"Email":         user.Email,
		}, true)
		if err != nil {
			//since we don't want to tell the user it failed, we'll log and move on
			logging.Error.Printf("Error sending email: %s\n", err)
		}
	}

	c.Status(http.StatusNoContent)
}

// @Summary Removes access to a server
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path string true "Server ID"
// @Param email path string true "Email of user"
// @Router /api/servers/{id}/users/{email} [delete]
func removeServerUser(c *gin.Context) {
	var err error
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}
	ps := &services.Permission{DB: db}

	email := c.Param("email")
	username := c.Param("username")
	if email == "" && username == "" {
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

	var user *models.User
	if email != "" {
		user, err = us.GetByEmail(email)
	} else {
		user, err = us.Get(username)
	}

	if err != nil && err == gorm.ErrRecordNotFound {
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

	es := services.GetEmailService()
	err = es.SendEmail(user.Email, "removedFromServer", map[string]interface{}{
		"Server": server,
	}, true)
	if err != nil {
		//since we don't want to tell the user it failed, we'll log and move on
		logging.Error.Printf("Error sending email: %s\n", err)
	}

	c.Status(http.StatusNoContent)
}

// @Summary Rename server
// @Description Renames a server
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path string false "Server ID"
// @Param name path string false "Server Name"
// @Router /api/servers/{id}/name [post]
func renameServer(c *gin.Context) {
	var err error

	t, exist := c.Get("server")
	if !exist {
		logging.Error.Printf("getting server for rename with err `%s`", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	server, ok := t.(*models.Server)
	if !ok {
		response.HandleError(c, pufferpanel.ErrUnknownError, http.StatusInternalServerError)
		return
	}

	name := c.Param("name")
	if name == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	t, exist = c.Get("db")
	if !exist {
		logging.Error.Printf("getting server for rename with err `%s`", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	db, ok := t.(*gorm.DB)
	if !ok {
		response.HandleError(c, pufferpanel.ErrUnknownError, http.StatusInternalServerError)
		return
	}
	ss := &services.Server{DB: db}

	server.Name = name
	err = ss.Update(server)
	if err != nil {
		logging.Error.Printf("renaming server with err `%s`", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

/*// @Summary Gets available OAuth2 scopes for the calling user
// @Description This allows a caller to see what scopes they have for a server, which can be used to generate a new OAuth2 client or just to know what they can do without making more calls
// @Accept json
// @Produce json
// @Success 200 {object} pufferpanel.Scopes
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path string true "Server ID"
// @Router /api/servers/{id}/oauth2 [get]
func getAvailableOauth2Perms(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	server := c.MustGet("server").(*models.Server)

	db := middleware.GetDatabase(c)
	ps := &services.Permission{DB: db}

	perms, err := ps.GetForUserAndServer(user.ID, &server.Identifier)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, perms.ToScopes())
}*/

// @Summary Gets server-level OAuth2 credentials for the logged in user
// @Accept json
// @Produce json
// @Success 200 {object} []models.Client
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path string true "Server ID"
// @Router /api/servers/{id}/oauth2 [get]
func getOAuth2Clients(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	server := c.MustGet("server").(*models.Server)

	db := middleware.GetDatabase(c)
	os := &services.OAuth2{DB: db}

	clients, err := os.GetForUserAndServer(user.ID, server.Identifier)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, clients)
}

// @Summary Creates server-level OAuth2 credentials for the logged in user
// @Accept json
// @Produce json
// @Success 200 {object} models.CreatedClient
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path string true "Server ID"
// @Param body body models.Client false "Client to create"
// @Router /api/servers/{id}/oauth2 [post]
func createOAuth2Client(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	server := c.MustGet("server").(*models.Server)

	db := middleware.GetDatabase(c)
	os := &services.OAuth2{DB: db}

	var request models.Client
	err := c.BindJSON(&request)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	client := &models.Client{
		ClientId: uuid.NewV4().String(),
		UserId:   user.ID,
		ServerId: sql.NullString{
			String: server.Identifier,
			Valid:  server.Identifier != "",
		},
		Name:        request.Name,
		Description: request.Description,
	}

	secret, err := pufferpanel.GenerateRandomString(36)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = client.SetClientSecret(secret)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = os.Update(client)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, models.CreatedClient{
		ClientId:     client.ClientId,
		ClientSecret: secret,
	})
}

// @Summary Deletes server-level OAuth2 credential
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path string true "Server ID"
// @Param clientId path string true "Client ID"
// @Router /api/servers/{id}/oauth2/{clientId} [delete]
func deleteOAuth2Client(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	server := c.MustGet("server").(*models.Server)
	clientId := c.Param("clientId")

	db := middleware.GetDatabase(c)
	os := &services.OAuth2{DB: db}

	cl, err := os.Get(clientId)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	//ensure the client id is specific for this server, and this user
	if cl.UserId != user.ID || !cl.ServerId.Valid || cl.ServerId.String != server.Identifier {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	err = os.Delete(cl)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
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
