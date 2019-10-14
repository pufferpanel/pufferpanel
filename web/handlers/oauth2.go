package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/v3"
	"github.com/pufferpanel/apufferi/v3/logging"
	"github.com/pufferpanel/apufferi/v3/response"
	"github.com/pufferpanel/apufferi/v3/scope"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/services"
	webHttp "net/http"
	"strconv"
	"strings"
)

func HasOAuth2Token(c *gin.Context) {
	//if there's a cookie with the token, use that
	cookie, _ := c.Cookie("puffer_auth")
	if cookie != "" {
		c.Request.Header.Set("Authorization", "Bearer "+cookie)
	}

	authHeader := c.Request.Header.Get("Authorization")
	authHeader = strings.TrimSpace(authHeader)

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		c.AbortWithStatus(403)
	}

	if parts[0] != "Bearer" || parts[1] == "" {
		c.AbortWithStatus(403)
	}

	token, err := services.ParseToken(parts[1])

	if err != nil || !token.Valid {
		c.AbortWithStatus(403)
		return
	}

	c.Set("token", token)
	c.Next()
}

func OAuth2Handler(requiredScope scope.Scope, requireServer bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		res := response.From(c)

		db, err := database.GetConnection()

		if response.HandleError(res, err) {
			return
		}

		token, ok := c.Get("token")
		if !ok {
			res.Status(webHttp.StatusInternalServerError).Message("error validating credentials").Fail()
			logging.Build(logging.ERROR).WithMessage("error validating credentials").WithError(err).Log()
			c.Abort()
			return
		}
		jwtToken, ok := token.(*apufferi.Token)
		if !ok {
			res.Status(webHttp.StatusInternalServerError).Message("error validating credentials").Fail()
			logging.Build(logging.ERROR).WithMessage("error validating credentials").WithError(err).Log()
			c.Abort()
			return
		}

		ti := jwtToken.Claims

		ss := &services.Server{DB: db}
		us := &services.User{DB: db}

		var serverId string

		var server *models.Server

		i := c.Param("serverId")
		if requireServer {
			server, err = ss.Get(i)
			if err != nil {
				res.Status(webHttp.StatusUnauthorized).Fail()
				c.Abort()
				return
			}
		}

		if requireServer && (server == nil || server.Identifier == "") {
			res.Status(webHttp.StatusUnauthorized).Fail()
			c.Abort()
			return
		} else if requireServer {
			serverId = server.Identifier
		}

		ui, err := strconv.Atoi(ti.Subject)
		if err != nil {
			res.Status(webHttp.StatusUnauthorized).Fail()
			c.Abort()
			return
		}

		user, err := us.GetById(uint(ui))
		if err != nil {
			res.Status(webHttp.StatusUnauthorized).Fail()
			c.Abort()
			return
		}

		allowed := false

		//if this is an audience of oauth2, we can use token directly
		if ti.Audience == "oauth2" {
			scopes := ti.PanelClaims.Scopes[serverId]
			if scopes != nil && apufferi.ContainsScope(scopes, requiredScope) {
				allowed = true
			} else {
				//if there isn't a defined rule, is this user an admin?
				scopes := ti.PanelClaims.Scopes[""]
				if scopes != nil && apufferi.ContainsScope(scopes, scope.ServersAdmin) {
					allowed = true
				}
			}
		} else if ti.Audience == "session" {
			//otherwise, we have to look at what the user has since session based
			ps := &services.Permission{DB: db}
			var perms *models.Permissions
			if serverId == "" {
				perms, err = ps.GetForUserAndServer(user.ID, nil)
			} else {
				perms, err = ps.GetForUserAndServer(user.ID, &serverId)
			}
			if err != nil {
				res.Status(webHttp.StatusInternalServerError).Message("error validating credentials").Fail()
				logging.Build(logging.ERROR).WithMessage("error validating credentials").WithError(err).Log()
				c.Abort()
				return
			}
			if apufferi.ContainsScope(perms.ToScopes(), requiredScope) {
				allowed = true
			} else {
				perms, err = ps.GetForUserAndServer(user.ID, nil)
				if err != nil {
					res.Status(webHttp.StatusInternalServerError).Message("error validating credentials").Fail()
					logging.Build(logging.ERROR).WithMessage("error validating credentials").WithError(err).Log()
					c.Abort()
					return
				}
				if apufferi.ContainsScope(perms.ToScopes(), scope.ServersAdmin) {
					allowed = true
				}
			}
		} else {
			res.Status(webHttp.StatusForbidden).Fail()
			c.Abort()
			return
		}

		if !allowed {
			res.Status(webHttp.StatusForbidden).Fail()
			c.Abort()
			return
		}

		c.Set("server", server)
		c.Set("user", user)
		c.Next()
	}
}
