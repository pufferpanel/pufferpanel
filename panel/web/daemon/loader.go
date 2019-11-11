package daemon

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/daemon/routing/server"
	"github.com/pufferpanel/pufferpanel/v2/panel/models"
	"github.com/pufferpanel/pufferpanel/v2/panel/services"
	"github.com/pufferpanel/pufferpanel/v2/panel/web/handlers"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"strings"
)

var rootEngine *gin.Engine

func RegisterRoutes(engine *gin.Engine, rg *gin.RouterGroup) {
	rootEngine = engine
	g := rg.Group("/server", handlers.HasOAuth2Token, handlers.NeedsDatabase)
	{
		//g.Any("", proxyServerRequest)
		g.Any("/:id", proxyServerRequest)
		g.Any("/:id/*path", proxyServerRequest)
	}

	g = rg.Group("/socket", handlers.HasOAuth2Token, handlers.NeedsDatabase)
	{
		//g.Any("", proxyServerRequest)
		g.Any("/:id", proxyServerRequest)
	}

	r := rg.Group("/node", handlers.HasOAuth2Token, handlers.NeedsDatabase)
	{
		//g.Any("", proxyServerRequest)
		r.Any("/:id", proxyNodeRequest)
		r.Any("/:id/*path", proxyNodeRequest)
	}

	if viper.GetBool("localNode") {
		l := rg.Group("/local")
		{
			server.RegisterRoutes(l)
		}
	}
}

func proxyServerRequest(c *gin.Context) {
	db := handlers.GetDatabase(c)
	ss := &services.Server{DB: db}
	ns := &services.Node{DB: db}
	ps := &services.Permission{DB: db}

	serverId := c.Param("id")
	if serverId == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	path := strings.TrimPrefix(c.Request.URL.Path, "/daemon")

	s, err := ss.Get(serverId)
	if err != nil && !gorm.IsRecordNotFoundError(err) && response.HandleError(c, err, http.StatusInternalServerError) {
		return
	} else if s == nil || s.Identifier == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	token := c.MustGet("token").(*pufferpanel.Token)

	//if a session-token, we need to convert it to an oauth2 token instead
	if token.Claims.Audience == "session" {
		newToken, err := ps.GenerateOAuthForUser(cast.ToUint(token.Claims.Subject), &s.Identifier)
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}

		//set new header
		c.Request.Header.Set("Authorization", "Bearer "+newToken)
	}

	if viper.GetBool("localNode") && s.Node.Local {
		c.Request.URL.Path = "/daemon/local/" + path
		rootEngine.HandleContext(c)
	} else {
		if c.GetHeader("Upgrade") == "websocket" {
			proxySocketRequest(c, path, ns, &s.Node)
		} else {
			proxyHttpRequest(c, path, ns, &s.Node)
		}
	}
}

func proxyNodeRequest(c *gin.Context) {
	path := c.Param("path")
	db := handlers.GetDatabase(c)
	ns := &services.Node{DB: db}
	ps := &services.Permission{DB: db}

	nodeId := c.Param("id")
	if nodeId == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	id, err := strconv.ParseUint(nodeId, 10, 32)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	node, err := ns.Get(uint(id))
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	token := c.MustGet("token").(*pufferpanel.Token)

	//if a session-token, we need to convert it to an oauth2 token instead
	if token.Claims.Audience == "session" {
		newToken, err := ps.GenerateOAuthForUser(cast.ToUint(token.Claims.Subject), nil)
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}

		//set new header
		c.Header("Authorization", "Bearer "+newToken)
	}

	if c.GetHeader("Upgrade") == "websocket" {
		proxySocketRequest(c, path, ns, node)
	} else {
		proxyHttpRequest(c, path, ns, node)
	}
}

func proxyHttpRequest(c *gin.Context, path string, ns *services.Node, node *models.Node) {
	callResponse, err := ns.CallNode(node, c.Request.Method, path, c.Request.Body, c.Request.Header)

	if response.HandleError(c, err, http.StatusInternalServerError) {
	}

	//Even though apache isn't going to be in place, we can't set certain headers
	newHeaders := make(map[string]string, 0)
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
}

func proxySocketRequest(c *gin.Context, path string, ns *services.Node, node *models.Node) {
	err := ns.OpenSocket(node, path, c.Writer, c.Request)
	response.HandleError(c, err, http.StatusInternalServerError)
}
