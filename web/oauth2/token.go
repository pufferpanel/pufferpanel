package oauth2

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/v3/response"
	"github.com/pufferpanel/apufferi/v3/scope"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/services"
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
	res := response.From(c)
	db, err := database.GetConnection()
	if response.HandleError(res, err) {
		return
	}

	os := services.GetOAuth(db)

	res.Discard()
	os.HandleHTTPTokenRequest(c.Writer, c.Request)
}

func handleInfoRequest(c *gin.Context) {
	res := response.From(c)
	db, err := database.GetConnection()
	if response.HandleError(res, err) {
		return
	}

	res.Discard()

	os := services.GetOAuth(db)

	if token := c.PostForm("token"); token != "" {
		info, _, _ := os.GetInfo(token)
		if info == nil {
			data := make(map[string]interface{})
			data["active"] = false
			c.JSON(200, data)
		}
		c.JSON(200, info)
	} else {
		c.Status(400)
		return
	}
}
