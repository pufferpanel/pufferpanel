package shared

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/http"
	netHttp "net/http"
	"strings"
)

func NotImplemented (c *gin.Context) {
	http.Respond(c).Fail().Status(netHttp.StatusNotImplemented).Message("not implemented")
}

func CreateOptions(options ...string) gin.HandlerFunc {

	replacement := make([]string, len(options) + 1)

	copy(replacement, options)

	replacement[len(options)] = "OPTIONS"
	response := strings.Join(replacement, ",")

	return func (c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", response)
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", response)
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(netHttp.StatusOK)
	}
}