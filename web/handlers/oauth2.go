package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/v3/logging"
	"github.com/pufferpanel/apufferi/v3/response"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/services"
	webHttp "net/http"
	"strings"
)

func OAuth2(scope string, requireServer bool) gin.HandlerFunc {
	return oauth2Handler(scope, requireServer, false)
}

func OAuth2WithLimit(scope string, requireServer bool) gin.HandlerFunc {
	return oauth2Handler(scope, requireServer, true)
}

func HasOAuth2Token(c *gin.Context) {
	//if there's a cookie with the token, use that
	cookie, _ := c.Cookie("puffer_auth")
	if cookie != "" {
		c.Request.Header.Set("Authorization", "Bearer "+cookie)
	}

	authHeader := c.Request.Header.Get("Authorization")
	authHeader = strings.TrimSpace(authHeader)

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") && strings.TrimPrefix(authHeader, "Bearer ") != "" {
		c.AbortWithStatus(403)
	} else {
		c.Next()
	}
}

func oauth2Handler(scope string, requireServer bool, permitWithLimit bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		res := response.From(c)

		db, err := database.GetConnection()

		if response.HandleError(res, err) {
			return
		}

		os := services.GetOAuth(db)
		ss := &services.Server{DB: db}

		var serverId *uint

		i := c.Param("serverId")
		server, exists, err := ss.Get(i)

		if server != nil {
			serverId = &server.ID
		}

		if requireServer && (server == nil || server.ID == 0 || !exists || err != nil) {
			res.Status(webHttp.StatusUnauthorized).Fail()
			c.Abort()
			return
		}

		ti, err := os.ValidationBearerToken(c.Request)
		if err != nil {
			res.Status(webHttp.StatusUnauthorized).Fail()
			c.Abort()
			return
		}

		ci, allowed, err := os.HasRights(ti.GetAccess(), serverId, scope)
		if err != nil {
			res.Status(webHttp.StatusInternalServerError).Message("error validating credentials").Fail()
			logging.Build(logging.ERROR).WithMessage("error validating credentials").WithError(err).Log()
			c.Abort()
			return
		}

		if !allowed && permitWithLimit {
			for _, v := range ci.ServerScopes {
				if v.Scope == scope {
					allowed = true
					break
				}
			}
		}

		if !allowed {
			if ci == nil {
				res.Status(webHttp.StatusUnauthorized).Fail()
				c.Abort()
			} else {
				res.Status(webHttp.StatusForbidden).Fail()
				c.Abort()
			}
		} else {
			c.Set("accessToken", ti.GetAccess())
			c.Set("server", server)
			c.Set("user", &ci.User)
			c.Next()
		}
	}
}
