package daemon

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/apufferi/v4"
	"github.com/pufferpanel/apufferi/v4/response"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"github.com/pufferpanel/pufferpanel/v2/web/handlers"
	"github.com/spf13/cast"
	"net/http"
	"strconv"
	"strings"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/server", handlers.HasOAuth2Token, handlers.NeedsDatabase)
	{
		//g.Any("", proxyServerRequest)
		g.Any("/:id", proxyServerRequest)
		g.Any("/:id/*path", proxyServerRequest)
	}
	r := rg.Group("/node", handlers.HasOAuth2Token, handlers.NeedsDatabase)
	{
		//g.Any("", proxyServerRequest)
		r.Any("/:id", proxyNodeRequest)
		r.Any("/:id/*path", proxyNodeRequest)
	}
}

func proxyServerRequest(c *gin.Context) {
	db := handlers.GetDatabase(c)
	ss := &services.Server{DB: db}
	ns := &services.Node{DB: db}

	serverId := c.Param("id")
	if serverId == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	path := "/server/" + serverId + c.Param("path")

	server, err := ss.Get(serverId)
	if err != nil && !gorm.IsRecordNotFoundError(err) && response.HandleError(c, err, http.StatusInternalServerError) {
		return
	} else if server == nil || server.Identifier == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	token := c.MustGet("token").(*apufferi.Token)

	//if a session-token, we need to convert it to an oauth2 token instead
	if token.Claims.Audience == "session" {
		newToken, err := services.GenerateOAuthForUser(cast.ToUint(token.Claims.Subject), &server.Identifier)
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}

		//set new header
		c.Request.Header.Set("Authorization", "Bearer "+newToken)
	}

	if c.GetHeader("Upgrade") == "websocket" {
		proxySocketRequest(c, path, ns, &server.Node)
	} else {
		proxyHttpRequest(c, path, ns, &server.Node)
	}
}

func proxyNodeRequest(c *gin.Context) {
	path := c.Param("path")
	db := handlers.GetDatabase(c)
	ns := &services.Node{DB: db}

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

	token := c.MustGet("token").(*apufferi.Token)

	//if a session-token, we need to convert it to an oauth2 token instead
	if token.Claims.Audience == "session" {
		newToken, err := services.GenerateOAuthForUser(cast.ToUint(token.Claims.Subject), nil)
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
