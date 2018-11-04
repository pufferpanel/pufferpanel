package oauth2

import (
	"github.com/gin-gonic/gin"
	builder "github.com/pufferpanel/apufferi/http"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/shared"
	"gopkg.in/oauth2.v3/server"
)

var oauth2Server *server.Server

func registerTokens(g *gin.RouterGroup) {
	g.POST("/token", handleTokenRequest)
	g.OPTIONS("/token", shared.CreateOptions("POST"))

	g.POST("/validate", shared.NotImplemented)
	g.OPTIONS("/validate", shared.CreateOptions("POST"))
}

func handleTokenRequest(c *gin.Context) {
	response := builder.Respond(c)

	var oauth services.OAuthService
	var err error
	if oauth, err = services.GetOAuthService(); shared.HandleError(response, err) {
		return
	}

	oauth.HandleHTTPTokenRequest(c.Writer, c.Request)
}