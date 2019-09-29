package oauth2

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/v3/response"
	"github.com/pufferpanel/apufferi/v3/scope"
	"github.com/pufferpanel/pufferpanel/v2/web/handlers"
)

func registerTokens(g *gin.RouterGroup) {
	g.POST("/token", handleTokenRequest)
	g.OPTIONS("/token", response.CreateOptions("POST"))

	g.POST("/validate", response.NotImplemented)
	g.OPTIONS("/validate", response.CreateOptions("POST"))

	g.POST("/info", handlers.OAuth2(scope.OAuth2Info, false), handleInfoRequest)
	g.OPTIONS("/info", response.CreateOptions("POST"))
}

func handleTokenRequest(c *gin.Context) {
}

func handleInfoRequest(c *gin.Context) {
}
