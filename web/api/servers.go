package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/database"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func registerServers(g *gin.RouterGroup) {
	g.Handle("GET", "", searchServers)
	g.Handle("OPTIONS", "", response.CreateOptions("GET"))

	g.Handle("GET", "/:serverId", middleware.RequiresPermission(pufferpanel.ScopeServerView), middleware.ResolveServerPanel, getServer)
	g.Handle("PUT", "/:serverId", middleware.RequiresPermission(pufferpanel.ScopeServerCreate), middleware.ResolveServerNode, middleware.HasTransaction, createServer)
	g.Handle("DELETE", "/:serverId", middleware.RequiresPermission(pufferpanel.ScopeServerDelete), middleware.ResolveServerPanel, middleware.HasTransaction, deleteServer)
	g.Handle("OPTIONS", "/:serverId", response.CreateOptions("PUT", "GET", "POST", "DELETE"))

	g.Handle("PUT", "/:serverId/name/:name", middleware.RequiresPermission(pufferpanel.ScopeServerEditName), middleware.ResolveServerPanel, middleware.HasTransaction, renameServer)
	g.Handle("OPTIONS", "/:serverId/name", response.CreateOptions("PUT"))
	g.Handle("OPTIONS", "/:serverId/name/:name", response.CreateOptions("PUT"))

	g.Handle("GET", "/:serverId/definition", middleware.RequiresPermission(pufferpanel.ScopeServerViewDefinition), middleware.ResolveServerPanel, proxyServerRequest)
	g.Handle("PUT", "/:serverId/definition", middleware.RequiresPermission(pufferpanel.ScopeServerEditDefinition), middleware.ResolveServerPanel, middleware.HasTransaction, editServer)
	g.Handle("OPTIONS", "/:serverId/definition", response.CreateOptions("PUT", "GET"))

	g.Handle("GET", "/:serverId/user", middleware.RequiresPermission(pufferpanel.ScopeServerUserView), middleware.ResolveServerPanel, getServerUsers)
	g.Handle("OPTIONS", "/:serverId/user", response.CreateOptions("GET"))

	g.Handle("GET", "/:serverId/user/:email", middleware.RequiresPermission(pufferpanel.ScopeServerUserView), middleware.ResolveServerPanel, getServerUsers)
	g.Handle("PUT", "/:serverId/user/:email", middleware.RequiresPermission(pufferpanel.ScopeServerUserEdit), middleware.ResolveServerPanel, middleware.HasTransaction, editServerUser)
	g.Handle("DELETE", "/:serverId/user/:email", middleware.RequiresPermission(pufferpanel.ScopeServerUserDelete), middleware.ResolveServerPanel, middleware.HasTransaction, removeServerUser)
	g.Handle("OPTIONS", "/:serverId/user/:email", response.CreateOptions("GET", "PUT", "DELETE"))

	g.Handle("GET", "/:serverId/oauth2", middleware.RequiresPermission(pufferpanel.ScopeServerClientView), middleware.ResolveServerPanel, getOAuth2Clients)
	g.Handle("POST", "/:serverId/oauth2", middleware.RequiresPermission(pufferpanel.ScopeServerClientEdit), middleware.ResolveServerPanel, middleware.HasTransaction, createOAuth2Client)
	g.Handle("OPTIONS", "/:serverId/oauth2", response.CreateOptions("GET", "POST"))

	g.Handle("DELETE", "/:serverId/oauth2/:clientId", middleware.RequiresPermission(pufferpanel.ScopeServerClientDelete), middleware.ResolveServerPanel, middleware.HasTransaction, deleteOAuth2Client)
	g.Handle("OPTIONS", "/:serverId/oauth2/:clientId", response.CreateOptions("DELETE"))

	g.GET("/:serverId/data", middleware.RequiresPermission(pufferpanel.ScopeServerViewData), middleware.ResolveServerPanel, proxyServerRequest)
	g.POST("/:serverId/data", middleware.RequiresPermission(pufferpanel.ScopeServerEditData), middleware.ResolveServerPanel, editServerData)
	g.OPTIONS("/:serverId/data", response.CreateOptions("GET", "POST"))

	g.GET("/:serverId/flags", middleware.RequiresPermission(pufferpanel.ScopeServerViewFlags), middleware.ResolveServerPanel, proxyServerRequest)
	g.POST("/:serverId/flags", middleware.RequiresPermission(pufferpanel.ScopeServerEditFlags), middleware.ResolveServerPanel, proxyServerRequest)
	g.OPTIONS("/:serverId/flags", response.CreateOptions("GET", "POST"))

	g.GET("/:serverId/tasks", middleware.RequiresPermission(pufferpanel.ScopeServerTaskView), middleware.ResolveServerPanel, proxyServerRequest)
	g.POST("/:serverId/tasks", middleware.RequiresPermission(pufferpanel.ScopeServerTaskCreate), middleware.ResolveServerPanel, proxyServerRequest)
	g.PUT("/:serverId/tasks/:taskId", middleware.RequiresPermission(pufferpanel.ScopeServerTaskEdit), middleware.ResolveServerPanel, proxyServerRequest)
	g.DELETE("/:serverId/tasks/:taskId", middleware.RequiresPermission(pufferpanel.ScopeServerTaskDelete), middleware.ResolveServerPanel, proxyServerRequest)
	g.OPTIONS("/:serverId/tasks", response.CreateOptions("GET", "POST", "PUT", "DELETE"))

	//g.POST("/:serverId/tasks/:taskId/run", middleware.OAuth2Handler(pufferpanel.ScopeServersEdit, true), RunServerTask)
	//g.OPTIONS("/:serverId/tasks/:taskId/run", response.CreateOptions("POST"))

	g.POST("/:serverId/reload", middleware.RequiresPermission(pufferpanel.ScopeServerReload), middleware.ResolveServerPanel, proxyServerRequest)
	g.OPTIONS("/:serverId/reload", response.CreateOptions("POST"))

	g.POST("/:serverId/start", middleware.RequiresPermission(pufferpanel.ScopeServerStart), middleware.ResolveServerPanel, proxyServerRequest)
	g.OPTIONS("/:serverId/start", response.CreateOptions("POST"))

	g.POST("/:serverId/stop", middleware.RequiresPermission(pufferpanel.ScopeServerStop), middleware.ResolveServerPanel, proxyServerRequest)
	g.OPTIONS("/:serverId/stop", response.CreateOptions("POST"))

	g.POST("/:serverId/kill", middleware.RequiresPermission(pufferpanel.ScopeServerKill), middleware.ResolveServerPanel, proxyServerRequest)
	g.OPTIONS("/:serverId/kill", response.CreateOptions("POST"))

	g.POST("/:serverId/install", middleware.RequiresPermission(pufferpanel.ScopeServerInstall), middleware.ResolveServerPanel, proxyServerRequest)
	g.OPTIONS("/:serverId/install", response.CreateOptions("POST"))

	g.GET("/:serverId/file/*filename", middleware.RequiresPermission(pufferpanel.ScopeServerFileView), middleware.ResolveServerPanel, proxyServerRequest)
	g.PUT("/:serverId/file/*filename", middleware.RequiresPermission(pufferpanel.ScopeServerFileEdit), middleware.ResolveServerPanel, proxyServerRequest)
	g.DELETE("/:serverId/file/*filename", middleware.RequiresPermission(pufferpanel.ScopeServerFileEdit), middleware.ResolveServerPanel, proxyServerRequest)
	g.POST("/:serverId/file/*filename", middleware.RequiresPermission(pufferpanel.ScopeServerFileEdit), middleware.ResolveServerPanel, proxyServerRequest)
	g.OPTIONS("/:serverId/file/*filename", response.CreateOptions("GET", "PUT", "DELETE", "POST"))

	g.GET("/:serverId/console", middleware.RequiresPermission(pufferpanel.ScopeServerConsole), middleware.ResolveServerPanel, proxyServerRequest)
	g.POST("/:serverId/console", middleware.RequiresPermission(pufferpanel.ScopeServerSendCommand), middleware.ResolveServerPanel, proxyServerRequest)
	g.OPTIONS("/:serverId/console", response.CreateOptions("GET", "POST"))

	g.GET("/:serverId/stats", middleware.RequiresPermission(pufferpanel.ScopeServerStats), middleware.ResolveServerPanel, proxyServerRequest)
	g.OPTIONS("/:serverId/stats", response.CreateOptions("GET"))

	g.GET("/:serverId/status", middleware.RequiresPermission(pufferpanel.ScopeServerStatus), middleware.ResolveServerPanel, proxyServerRequest)
	g.OPTIONS("/:serverId/status", response.CreateOptions("GET"))

	g.POST("/:serverId/archive/*filename", middleware.RequiresPermission(pufferpanel.ScopeServerFileEdit), middleware.ResolveServerPanel, proxyServerRequest)
	g.OPTIONS("/:serverId/archive/*filename", response.CreateOptions("POST"))

	g.POST("/:serverId/extract/*filename", middleware.RequiresPermission(pufferpanel.ScopeServerFileEdit), middleware.ResolveServerPanel, proxyServerRequest)
	g.OPTIONS("/:serverId/extract/*filename", response.CreateOptions("POST"))

	p := g.Group("/:serverId/socket")
	{
		p.GET("", middleware.RequiresPermission(pufferpanel.ScopeServerView), cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowCredentials: true,
		}), middleware.ResolveServerPanel, proxyServerRequest)
		p.Handle("CONNECT", "", middleware.RequiresPermission(pufferpanel.ScopeServerView), func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Credentials", "false")
		})
		p.OPTIONS("", response.CreateOptions("GET"))
	}
}

// @Summary Search servers
// @Description Gets servers, and allowing for filtering of servers. * is a wildcard that can be used for text inputs
// @Success 200 {object} models.ServerSearchResponse
// @Param username query string false "Username to filter on, default is current user if NOT admin"
// @Param node query uint false "Node ID to filter on"
// @Param name query string false "Name of server to filter on"
// @Param limit query uint false "Max number of results to return"
// @Param page query uint false "What page to get back for many results"
// @Router /api/servers [get]
// @Security OAuth2Application[server.view]
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
		if pufferpanel.ContainsScope(p.Scopes, pufferpanel.ScopeAdmin) {
			isAdmin = true
		}
	}

	if !isAdmin && username != "" && user.Username != username {
		c.JSON(http.StatusOK, &models.ServerSearchResponse{
			Servers: []*models.ServerView{},
			Metadata: &pufferpanel.Metadata{Paging: &pufferpanel.Paging{
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
		Metadata: &pufferpanel.Metadata{Paging: &pufferpanel.Paging{
			Page:    uint(page),
			Size:    uint(pageSize),
			MaxSize: MaxPageSize,
			Total:   total,
		}},
	})
}

// @Summary Get a server
// @Description Gets a particular server
// @Success 200 {object} models.GetServerResponse
// @Param id path string true "Server ID"
// @Router /api/servers/{id} [get]
// @Security OAuth2Application[server.view]
func getServer(c *gin.Context) {
	server := getServerFromGin(c)

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

	c.JSON(http.StatusOK, d)
}

// @Summary Create server
// @Description Creates a server
// @Success 200 {object} models.CreateServerResponse
// @Param id path string true "Server ID"
// @Param server body models.ServerCreation true "Creation information"
// @Router /api/servers/{id} [put]
// @Security OAuth2Application[server.create]
func createServer(c *gin.Context) {
	var err error
	db := middleware.GetDatabase(c)
	ss := &services.Server{DB: db}
	ns := &services.Node{DB: db}
	us := &services.User{DB: db}
	ps := &services.Permission{DB: db}

	serverId := c.Param("serverId")

	if serverId == "" {
		gen, err := uuid.NewV4()
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}
		serverId = gen.String()[:8]
	}

	postBody := &models.ServerCreation{}
	err = c.Bind(postBody)
	postBody.Identifier = serverId
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	node, err := ns.Get(postBody.NodeId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
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
		IP:         cast.ToString(ip),
		Port:       cast.ToUint16(port),
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

		perm.Scopes = []*pufferpanel.Scope{
			pufferpanel.ScopeServerView,
			pufferpanel.ScopeServerViewData,
			pufferpanel.ScopeServerEditData,
			pufferpanel.ScopeServerEditFlags,
			pufferpanel.ScopeServerEditName,
			pufferpanel.ScopeServerViewData,
			pufferpanel.ScopeServerClientView,
			pufferpanel.ScopeServerClientEdit,
			pufferpanel.ScopeServerClientCreate,
			pufferpanel.ScopeServerClientDelete,
			pufferpanel.ScopeServerUserView,
			pufferpanel.ScopeServerUserCreate,
			pufferpanel.ScopeServerUserEdit,
			pufferpanel.ScopeServerUserDelete,
			pufferpanel.ScopeServerTaskView,
			pufferpanel.ScopeServerTaskRun,
			pufferpanel.ScopeServerTaskCreate,
			pufferpanel.ScopeServerTaskDelete,
			pufferpanel.ScopeServerReload,
			pufferpanel.ScopeServerStart,
			pufferpanel.ScopeServerStop,
			pufferpanel.ScopeServerKill,
			pufferpanel.ScopeServerInstall,
			pufferpanel.ScopeServerFileView,
			pufferpanel.ScopeServerFileEdit,
			pufferpanel.ScopeServerSftp,
			pufferpanel.ScopeServerConsole,
			pufferpanel.ScopeServerSendCommand,
			pufferpanel.ScopeServerStats,
			pufferpanel.ScopeServerStatus,
		}

		err = ps.UpdatePermissions(perm)
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}
	}

	err = db.Commit().Error
	if response.HandleError(c, err, http.StatusInternalServerError) {
		c.Abort()
		return
	}

	data, _ := json.Marshal(postBody.Server)
	reader := io.NopCloser(bytes.NewReader(data))

	nodeResponse, err := ns.CallNode(node, "PUT", "/daemon/server/"+server.Identifier, reader, c.Request.Header)
	if nodeResponse != nil {
		defer pufferpanel.Close(nodeResponse.Body)
	}

	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	if nodeResponse.StatusCode != http.StatusOK {
		resData, err := io.ReadAll(nodeResponse.Body)
		if err != nil {
			logging.Error.Printf("Failed to parse response from daemon\n%s", err.Error())
		}
		logging.Error.Printf("Unexpected response from daemon: %+v\n%s", nodeResponse.StatusCode, string(resData))
		//assume daemon gives us a valid response, directly forward to client
		c.Header("Content-Type", "application/json")
		c.Status(nodeResponse.StatusCode)
		_, _ = c.Writer.Write(resData)
		c.Abort()

		//revert DB, we need to actually now use the non-transactional connection
		noTransDbCtx, _ := c.Get("noTransactionDb")
		noTransDb := noTransDbCtx.(*gorm.DB)

		tempSS := &services.Server{DB: noTransDb}
		_ = tempSS.Delete(server.Identifier)
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

// @Summary Update server definition
// @Description Updates a server definition
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Param server body models.ServerWithName true "Server definition"
// @Router /api/servers/{id}/definition [put]
// @Security OAuth2Application[server.edit.definition]
func editServer(c *gin.Context) {
	var err error
	db := middleware.GetDatabase(c)
	ss := &services.Server{DB: db}
	ns := &services.Node{DB: db}

	server := getServerFromGin(c)

	postBody := &models.ServerWithName{}
	err = c.Bind(postBody)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	postBody.Identifier = server.Identifier

	port, err := getFromDataOrDefault(postBody.Variables, "port", uint16(0))
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}
	server.Port = cast.ToUint16(port)

	ip, err := getFromDataOrDefault(postBody.Variables, "ip", "0.0.0.0")
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}
	server.IP = cast.ToString(ip)

	if postBody.Name != "" {
		server.Name = postBody.Name
	}

	err = ss.Update(server)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	data, _ := json.Marshal(postBody)
	reader := io.NopCloser(bytes.NewReader(data))

	nodeResponse, err := ns.CallNode(&server.Node, "PUT", "/daemon/server/"+postBody.Identifier+"/definition", reader, c.Request.Header)
	defer pufferpanel.CloseResponse(nodeResponse)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	if nodeResponse.StatusCode != http.StatusOK {
		resData, err := io.ReadAll(nodeResponse.Body)
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
	c.Status(http.StatusNoContent)
}

// @Summary Deletes a server
// @Description Deletes a server from the panel
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Router /api/servers/{id} [delete]
// @Security OAuth2Application[server.delete]
func deleteServer(c *gin.Context) {
	var err error

	db := middleware.GetDatabase(c)
	ss := &services.Server{DB: db}
	ns := &services.Node{DB: db}

	server := getServerFromGin(c)

	node, err := ns.Get(server.NodeID)
	if err != nil {
		response.HandleError(c, err, http.StatusInternalServerError)
		return
	}

	//we need to know what users are impacted by a server being deleted
	ps := services.Permission{DB: db}
	users := make([]models.User, 0)
	perms, err := ps.GetForServer(server.Identifier)
	if err != nil {
		response.HandleError(c, err, http.StatusInternalServerError)
		return
	}
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

	_, skipNode := c.GetQuery("skipNode")
	if !skipNode {
		newHeader, err := c.Cookie("puffer_auth")
		if response.HandleError(c, err, http.StatusInternalServerError) {
			db.Rollback()
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
	}

	err = ss.Delete(server.Identifier)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		db.Rollback()
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
// @Success 200 {object} []models.PermissionView
// @Param id path string true "Server ID"
// @Router /api/servers/{id}/user [get]
// @Security OAuth2Application[server.users.view]
func getServerUsers(c *gin.Context) {
	var err error
	db := middleware.GetDatabase(c)
	ps := &services.Permission{DB: db}

	server := getServerFromGin(c)

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
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Param email path string true "Email of user"
// @Param body body models.PermissionView true "New permissions to apply"
// @Router /api/servers/{id}/users/{email} [put]
// @Security OAuth2Application[server.users.edit]
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

	server := getServerFromGin(c)

	currentUser := c.MustGet("user").(*models.User)
	currentPerms, err := ps.GetForUserAndServer(currentUser.ID, &server.Identifier)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	var registerToken string
	var user *models.User
	if email != "" {
		user, err = us.GetByEmail(email)
	} else {
		user, err = us.Get(username)
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) && response.HandleError(c, err, http.StatusInternalServerError) {
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		if email == "" {
			response.HandleError(c, err, http.StatusBadRequest)
			return
		}
		//we need to create the user here, since it's a new email we've not seen

		un, err := uuid.NewV4()
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}
		user = &models.User{
			Username: un.String(),
			Email:    email,
		}
		token, err := uuid.NewV4()
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}
		registerToken = token.String()
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

	//update perms to match this "setup", but not stomp over what the user can't change
	replacement := pufferpanel.UpdateScopesWhereGranted(existing.Scopes, perms.Scopes, currentPerms.Scopes)

	existing.Scopes = replacement
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
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Param email path string true "Email of user"
// @Router /api/servers/{id}/users/{email} [delete]
// @Security OAuth2Application[server.users.delete]
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

	server := getServerFromGin(c)

	var user *models.User
	if email != "" {
		user, err = us.GetByEmail(email)
	} else {
		user, err = us.Get(username)
	}

	if err != nil && response.HandleError(c, err, http.StatusInternalServerError) {
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
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Param name body pufferpanel.Name true "New server name"
// @Router /api/servers/{id}/name [post]
// @Security OAuth2Application[server.edit.name]
func renameServer(c *gin.Context) {
	var err error

	server := getServerFromGin(c)

	var name pufferpanel.Name
	err = c.ShouldBindJSON(&name)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	t, exist := c.Get("db")
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

	server.Name = name.Name
	err = ss.Update(server)
	if err != nil {
		logging.Error.Printf("renaming server with err `%s`", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Gets server-level OAuth2 credentials for the logged in user
// @Success 200 {object} []models.Client
// @Param id path string true "Server ID"
// @Router /api/servers/{id}/oauth2 [get]
// @Security OAuth2Application[server.clients.view]
func getOAuth2Clients(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	server := getServerFromGin(c)

	db := middleware.GetDatabase(c)
	os := &services.OAuth2{DB: db}

	clients, err := os.GetForUserAndServer(user.ID, server.Identifier)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, clients)
}

// @Summary Creates server-level OAuth2 credentials for the logged in user
// @Success 200 {object} models.CreatedClient
// @Param id path string true "Server ID"
// @Param body body models.Client false "Client to create"
// @Router /api/servers/{id}/oauth2 [post]
// @Security OAuth2Application[server.clients.create]
func createOAuth2Client(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	server := getServerFromGin(c)

	db := middleware.GetDatabase(c)
	os := &services.OAuth2{DB: db}

	var request models.Client
	err := c.BindJSON(&request)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	id, err := uuid.NewV4()
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}
	client := &models.Client{
		ClientId:    id.String(),
		UserId:      user.ID,
		ServerId:    server.Identifier,
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
// @Success 204 {object} nil
// @Param id path string true "Server ID"
// @Param clientId path string true "Client ID"
// @Router /api/servers/{id}/oauth2/{clientId} [delete]
// @Security OAuth2Application[server.clients.delete]
func deleteOAuth2Client(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	server := getServerFromGin(c)
	clientId := c.Param("clientId")

	db := middleware.GetDatabase(c)
	os := &services.OAuth2{DB: db}

	client, err := os.Get(clientId)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	//ensure the client id is specific for this server, and this user
	if client.UserId != user.ID || client.ServerId != server.Identifier {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	err = os.Delete(client)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Update server data
// @Description Updates a server's set of variables
// @Success 202 {object} nil
// @Param id path string true "Server ID"
// @Param server body pufferpanel.ServerData true "Server variables"
// @Router /api/servers/{id}/data [put]
// @Security OAuth2Application[server.edit.data]
func editServerData(c *gin.Context) {
	server := getServerFromGin(c)

	//clone request body, so we can re-set it for the proxy call
	useHere := &bytes.Buffer{}
	useThere := &bytes.Buffer{}

	multi := io.MultiWriter(useHere, useThere)
	_, err := io.CopyN(multi, c.Request.Body, 1024 /* 1KB */ *512 /* .5 MB */)
	if err != nil && errors.Is(err, io.EOF) && response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	_ = c.Request.Body.Close()
	c.Request.Body = io.NopCloser(useThere)

	postBody := &pufferpanel.ServerData{}
	err = json.NewDecoder(useHere).Decode(postBody)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	port, err := getFromDataOrDefault(postBody.Variables, "port", uint16(0))
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}
	server.Port = cast.ToUint16(port)

	ip, err := getFromDataOrDefault(postBody.Variables, "ip", "0.0.0.0")
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}
	server.IP = cast.ToString(ip)

	db := middleware.GetDatabase(c)
	ss := &services.Server{DB: db}
	err = ss.Update(server)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	proxyServerRequest(c)
}

func getFromData(variables map[string]pufferpanel.Variable, key string) (result interface{}, exists bool) {
	for k, v := range variables {
		if k == key {
			return v.Value, true
		}
	}
	return nil, false
}

func getFromDataOrDefault(variables map[string]pufferpanel.Variable, key string, val interface{}) (interface{}, error) {
	res, exists := getFromData(variables, key)

	if exists {
		return pufferpanel.Convert(res, val)
	}

	return val, nil
}

func proxyServerRequest(c *gin.Context) {
	db := middleware.GetDatabase(c)
	ns := &services.Node{DB: db}

	resolvedPath := "/daemon/server/" + strings.TrimPrefix(c.Request.RequestURI, "/api/servers/")

	server := c.MustGet("server").(*models.Server)
	node := &server.Node

	ps := &services.PanelService{}
	token := ps.GetActiveToken()

	//switch to our token for auth
	c.Request.Header.Set("Authorization", "Bearer "+token)

	if node.IsLocal() {
		c.Request.URL.Path = resolvedPath
		pufferpanel.Engine.HandleContext(c)
	} else {
		if c.IsWebsocket() {
			//for websocket, nuke the query params to avoid trying to escalate
			resolvedPath = strings.SplitN(c.Request.URL.Path, "?", 2)[0]
			if !strings.HasPrefix(resolvedPath, "/") {
				resolvedPath = "/" + resolvedPath
			}

			//add the params we can grant for this request
			scopes := c.MustGet("scopes").([]*pufferpanel.Scope)
			var params []string
			if pufferpanel.ContainsScope(scopes, pufferpanel.ScopeServerConsole) {
				params = append(params, "console")
			}
			if pufferpanel.ContainsScope(scopes, pufferpanel.ScopeServerStatus) {
				params = append(params, "status")
			}
			if pufferpanel.ContainsScope(scopes, pufferpanel.ScopeServerStats) {
				params = append(params, "stats")
			}
			resolvedPath = resolvedPath + "?" + strings.Join(params, "&")

			proxySocketRequest(c, resolvedPath, ns, node)
		} else {
			proxyHttpRequest(c, resolvedPath, ns, node)
		}
	}
	c.Abort()
}

func proxyHttpRequest(c *gin.Context, path string, ns *services.Node, node *models.Node) {
	callResponse, err := ns.CallNode(node, c.Request.Method, path, c.Request.Body, c.Request.Header)

	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	defer pufferpanel.CloseResponse(callResponse)

	//Even though apache isn't going to be in place, we can't set certain headers
	newHeaders := make(map[string]string)
	for k, v := range callResponse.Header {
		switch k {
		case "Transfer-Encoding":
		case "Content-Type":
		case "Content-Length":
			continue
		default:
			newHeaders[k] = strings.Join(v, ", ")
		}
	}

	c.DataFromReader(callResponse.StatusCode, callResponse.ContentLength, callResponse.Header.Get("Content-Type"), callResponse.Body, newHeaders)
	c.Abort()
}

func proxySocketRequest(c *gin.Context, path string, ns *services.Node, node *models.Node) {
	if node.IsLocal() {
		//have gin handle the request again, but send it to daemon instead
		c.Request.URL.Path = path
		pufferpanel.Engine.HandleContext(c)
	} else {
		err := ns.OpenSocket(node, path, c.Writer, c.Request)
		response.HandleError(c, err, http.StatusInternalServerError)
	}
	c.Abort()
}

func getServerFromGin(c *gin.Context) *models.Server {
	return c.MustGet("server").(*models.Server)
}
