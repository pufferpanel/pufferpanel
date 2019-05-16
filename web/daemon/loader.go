package daemon

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/apufferi/response"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/shared"
	netHttp "net/http"
	"strings"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/server")
	{
		//g.Any("", proxyServerRequest)
		g.Any("/:id", proxyServerRequest)
		g.Any("/:id/*path", proxyServerRequest)
	}
}

func proxyServerRequest(c *gin.Context) {
	res := response.Respond(c)

	serverId := c.Param("id")
	if serverId == "" {
		c.AbortWithStatus(404)
		return
	}

	path := "/server/" + serverId + c.Param("path")

	ss, err := services.GetServerService()
	if shared.HandleError(res, err) {
		return
	}

	ns, err := services.GetNodeService()
	if shared.HandleError(res, err) {
		return
	}

	server, exists, err := ss.Get(serverId)
	if err != nil && !gorm.IsRecordNotFoundError(err) && shared.HandleError(res, err) {
		return
	} else if !exists || server == nil {
		response.Respond(c).Status(netHttp.StatusNotFound).Fail().Send()
		return
	}

	callResponse, err := ns.CallNode(&server.Node, c.Request.Method, path, c.Request.Body, c.Request.Header)

	//this only will throw an error if we can't get to the node
	//so if error, use our response messenger, otherwise copy response from node to client
	if err != nil {
		response.Respond(c).Status(netHttp.StatusInternalServerError).Fail().Message(err.Error()).Send()
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

	c.DataFromReader(callResponse.StatusCode, callResponse.ContentLength, callResponse.Header.Get("Content-Type"), callResponse.Body, newHeaders)
}
