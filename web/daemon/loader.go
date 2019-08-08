package daemon

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/apufferi/response"
	"github.com/pufferpanel/pufferpanel/errors"
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/shared"
	netHttp "net/http"
	"strconv"
	"strings"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/server")
	{
		//g.Any("", proxyServerRequest)
		g.Any("/:id", proxyServerRequest)
		g.Any("/:id/*path", proxyServerRequest)
	}
	r := rg.Group("/node")
	{
		//g.Any("", proxyServerRequest)
		r.Any("/:id", proxyNodeRequest)
		r.Any("/:id/*path", proxyNodeRequest)
	}
}

func proxyServerRequest(c *gin.Context) {
	res := response.From(c)

	serverId := c.Param("id")
	if serverId == "" {
		res.Fail().Status(404)
		return
	}

	path := "/server/" + serverId + c.Param("path")

	ss, err := services.GetServerService()
	if ss == nil && err == nil {
		err = errors.ErrServiceNotAvailable
	}
	if shared.HandleError(res, err) {
		return
	}

	ns, err := services.GetNodeService()
	if ns == nil && err == nil {
		err = errors.ErrServiceNotAvailable
	}
	if shared.HandleError(res, err) {
		return
	}

	server, exists, err := ss.Get(serverId)
	if err != nil && !gorm.IsRecordNotFoundError(err) && shared.HandleError(res, err) {
		return
	} else if !exists || server == nil {
		res.Status(netHttp.StatusNotFound).Fail()
		return
	}

	if c.GetHeader("Upgrade") == "websocket" {
		proxySocketRequest(c, path, ns, &server.Node)
	} else {
		proxyHttpRequest(c, path, ns, &server.Node)
	}
}

func proxyNodeRequest(c *gin.Context) {
	path := c.Param("path")

	res := response.From(c)

	nodeId := c.Param("id")
	if nodeId == "" {
		res.Status(404).Fail()
		return
	}

	ns, err := services.GetNodeService()
	if ns == nil && err == nil {
		err = errors.ErrServiceNotAvailable
	}
	if shared.HandleError(res, err) {
		return
	}

	id, err := strconv.ParseUint(nodeId, 10, 32)
	if shared.HandleError(res, err) {
		return
	}

	node, exists, err := ns.Get(uint(id))
	if shared.HandleError(res, err) {
		return
	} else if !exists {
		res.Fail().Status(404)
		return
	}

	if c.GetHeader("Upgrade") == "websocket" {
		proxySocketRequest(c, path, ns, node)
	} else {
		proxyHttpRequest(c, path, ns, node)
	}
}

func proxyHttpRequest(c *gin.Context, path string, ns services.NodeService, node *models.Node) {
	callResponse, err := ns.CallNode(node, c.Request.Method, path, c.Request.Body, c.Request.Header)

	//this only will throw an error if we can't get to the node
	//so if error, use our response messenger, otherwise copy response from node to client
	if err != nil {
		response.From(c).Status(netHttp.StatusInternalServerError).Fail().Error(err)
		return
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

	response.From(c).Discard()
	c.DataFromReader(callResponse.StatusCode, callResponse.ContentLength, callResponse.Header.Get("Content-Type"), callResponse.Body, newHeaders)
}

func proxySocketRequest(c *gin.Context, path string, ns services.NodeService, node *models.Node) {
	err := ns.OpenSocket(node, path, c.Writer, c.Request)
	if err != nil {
		logging.Exception("error opening socket", err)
		response.From(c).Status(netHttp.StatusInternalServerError).Fail().Error(err)
		return
	}
	response.From(c).Discard()
}
