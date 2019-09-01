package oauth2

import (
	"github.com/gin-gonic/gin"
	builder "github.com/pufferpanel/apufferi/response"
	"github.com/pufferpanel/pufferpanel"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/web/handlers"
)

func registerTokens(g *gin.RouterGroup) {
	g.POST("/token", handleTokenRequest)
	g.OPTIONS("/token", pufferpanel.CreateOptions("POST"))

	g.POST("/validate", pufferpanel.NotImplemented)
	g.OPTIONS("/validate", pufferpanel.CreateOptions("POST"))

	g.POST("/info", handlers.OAuth2("oauth2.info", false), handleInfoRequest)
	g.OPTIONS("/info", pufferpanel.CreateOptions("POST"))
}

func handleTokenRequest(c *gin.Context) {
	response := builder.From(c)
	db, err := database.GetConnection()
	if pufferpanel.HandleError(response, err) {
		return
	}

	os := services.GetOAuth(db)

	response.Discard()
	os.HandleHTTPTokenRequest(c.Writer, c.Request)
}

func handleInfoRequest(c *gin.Context) {
	response := builder.From(c)
	db, err := database.GetConnection()
	if pufferpanel.HandleError(response, err) {
		return
	}

	response.Discard()

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
