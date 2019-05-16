package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/apufferi/response"
	"github.com/pufferpanel/pufferpanel/errors"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/shared"
	webHttp "net/http"
	"strconv"
	"strings"
)

func OAuth2(scope string, requireServer bool) gin.HandlerFunc {
	return oauth2Handler(scope, requireServer, false)
}

func OAuth2WithLimit(scope string, requireServer bool) gin.HandlerFunc {
	return oauth2Handler(scope, requireServer, true)
}

func HasOAuth2Token(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") && strings.TrimPrefix(authHeader, "Bearer ") != "" {
		c.AbortWithStatus(403)
	} else {
		c.Next()
	}
}

func oauth2Handler (scope string, requireServer bool, permitWithLimit bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		os, err := services.GetOAuthService()
		if err != nil || os == nil {
			if err == nil {
				err = errors.ErrServiceNotAvailable
			}

			shared.HandleError(response.Respond(c), err)
			c.Abort()
			return
		}

		var serverId *uint
		var id uint

		i := c.Query("serverId")
		if i != "" {
			t, _ := strconv.Atoi(i)
			id = uint(t)
			serverId = &id
		}

		if requireServer && serverId == nil {
			response.Respond(c).Status(webHttp.StatusUnauthorized).Fail().Send()
			c.Abort()
			return
		}

		ti, err := os.ValidationBearerToken(c.Request)
		if err != nil {
			response.Respond(c).Status(webHttp.StatusUnauthorized).Fail().Send()
			c.Abort()
			return
		}

		ci, allowed, err := os.HasRights(ti.GetAccess(), serverId, scope)
		if err != nil {
			response.Respond(c).Status(webHttp.StatusInternalServerError).Message("error validating credentials").Fail().Send()
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
				response.Respond(c).Status(webHttp.StatusUnauthorized).Fail().Send()
				c.Abort()
			} else {
				response.Respond(c).Status(webHttp.StatusForbidden).Fail().Send()
				c.Abort()
			}
		} else {
			c.Set("accessToken", ti.GetAccess())
			c.Next()
		}
	}
}