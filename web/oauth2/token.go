package oauth2

import (
	"github.com/gin-gonic/gin"
	builder "github.com/pufferpanel/apufferi/http"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/shared"
	"github.com/pufferpanel/pufferpanel/web/handlers"
)

func registerTokens(g *gin.RouterGroup) {
	g.POST("/token", handleTokenRequest)
	g.OPTIONS("/token", shared.CreateOptions("POST"))

	g.POST("/validate", shared.NotImplemented)
	g.OPTIONS("/validate", shared.CreateOptions("POST"))

	g.POST("/info", handlers.OAuth2("oauth2.info", false), handleInfoRequest)
	g.OPTIONS("/info", shared.CreateOptions("POST"))
}

func handleTokenRequest(c *gin.Context) {
	var oauth services.OAuthService
	var err error
	if oauth, err = services.GetOAuthService(); shared.HandleError(builder.Respond(c), err) {
		return
	}

	oauth.HandleHTTPTokenRequest(c.Writer, c.Request)
}

func handleInfoRequest(c *gin.Context) {
	var oauth services.OAuthService
	var err error
	if oauth, err = services.GetOAuthService(); err != nil {
		c.Status(500)
		return
	}

	if token := c.PostForm("token"); token != "" {
		info, _, _ := oauth.GetInfo(token)
		c.JSON(200, info)
	} else {
		c.Status(400)
		return
	}
}