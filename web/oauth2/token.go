package oauth2

import (
	"github.com/gin-gonic/gin"
	builder "github.com/pufferpanel/apufferi/http"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/shared"
	"net/http"
)

func registerTokens(g *gin.RouterGroup) {
	g.POST("/token", handleTokenRequest)
	g.OPTIONS("/token", shared.CreateOptions("POST"))

	g.POST("/validate", shared.NotImplemented)
	g.OPTIONS("/validate", shared.CreateOptions("POST"))

	g.POST("/info", handleInfoRequest)
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
	response := builder.Respond(c)

	var oauth services.OAuthService
	var err error
	if oauth, err = services.GetOAuthService(); shared.HandleError(response, err) {
		return
	}

	type body struct {
		token string
	}
	msg := &body{}

	if token := c.PostForm("token"); token != "" {
		msg.token = token
	} else {
		err := c.BindJSON(msg)
		if err != nil {
			response.Status(http.StatusBadRequest).Fail().Message(err.Error()).Send()
			return
		}
	}

	info, _, err := oauth.GetInfo(msg.token)
	if err != nil {
		response.Status(http.StatusBadRequest).Fail().Message(err.Error()).Send()
		return
	}

	response.Data(info).Send()
}