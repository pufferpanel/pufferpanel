package oauth2

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/v3/response"
)

func registerTokens(g *gin.RouterGroup) {
	g.POST("/token", handleTokenRequest)
	g.OPTIONS("/token", response.CreateOptions("POST"))
}

func handleTokenRequest(c *gin.Context) {
}
