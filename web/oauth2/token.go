package oauth2

import (
	"github.com/gin-gonic/gin"
	builder "github.com/pufferpanel/apufferi/http"
	"github.com/pufferpanel/pufferpanel/shared"
)

func registerTokens(g *gin.RouterGroup) {
	g.POST("/request", requestToken)
	g.OPTIONS("/request", shared.CreateOptions("POST"))

	g.POST("/info", tokenInfo)
	g.OPTIONS("/info", shared.CreateOptions("POST"))
}

func requestToken(c *gin.Context) {
	response := builder.Respond(c)

	response.Send()
}


func tokenInfo(c *gin.Context) {
	response := builder.Respond(c)

	response.Send()
}