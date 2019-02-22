package daemon

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/apufferi/http"
	"github.com/pufferpanel/pufferpanel/services"
	netHttp "net/http"
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
	serverId := c.Param("id")
	if serverId == "" {
		c.AbortWithStatus(404)
		return
	}

	ss, err := services.GetServerService()
	if err != nil {
		http.Respond(c).Status(netHttp.StatusInternalServerError).Fail().Message(err.Error()).Send()
	}

	_, exists, err := ss.Get(serverId)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		http.Respond(c).Status(netHttp.StatusInternalServerError).Fail().Message(err.Error()).Send()
	} else if !exists {
		http.Respond(c).Status(netHttp.StatusNotFound).Fail().Send()
	}
}