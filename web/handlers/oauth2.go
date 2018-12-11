package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/http"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferpanel/services"
	webHttp "net/http"
	"strconv"
)

func OAuth2(scope string, requireServer bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		os, err := services.GetOAuthService()
		if err != nil {
			http.Respond(c).Status(webHttp.StatusInternalServerError).Message("oauth2 service is not available").Fail().Send()
			logging.Errorf("oauth2 service is not available", err)
			c.Abort()
		}

		var serverId *uint
		var id uint

		if requireServer {
			i := c.Query("id")
			t, _ := strconv.Atoi(i)
			id = uint(t)
			serverId = &id
		}

		ci, allowed, err := os.HasRights("", serverId, scope)
		if err != nil {
			http.Respond(c).Status(webHttp.StatusInternalServerError).Message("error validating credentials").Fail().Send()
			logging.Errorf("error validating credentials", err)
			c.Abort()
		}

		if !allowed {
			if ci == nil {
				http.Respond(c).Status(webHttp.StatusUnauthorized).Fail().Send()
				c.Abort()
			} else {
				http.Respond(c).Status(webHttp.StatusForbidden).Fail().Send()
				c.Abort()
			}
		} else {
			c.Next()
		}
	}
}