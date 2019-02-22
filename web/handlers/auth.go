package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferpanel/services"
	"strings"
	"time"
)

var noLogin = []string{"/auth/","/error/", "/daemon/", "/api/"}
var assetFiles = []string{".js", ".css", ".img", ".ico", ".png", ".gif"}

func AuthMiddleware(c *gin.Context) {
	for _, v := range noLogin {
		if strings.HasPrefix(c.Request.URL.Path, v) {
			c.Next()
			return
		}
	}

	cookie, err := c.Cookie("puffer_auth")

	if err != nil || cookie == "" {
		//determine if it's an asset, otherwise, we can redirect if it's a GET
		//dev only requirement?
		for _, v := range assetFiles {
			if strings.HasSuffix(c.Request.URL.Path, v) {
				c.AbortWithStatus(404)
				return
			}
		}

		if c.Request.Method != "GET" {
			c.AbortWithStatus(403)
 		} else {
 			c.Redirect(302, "/auth/login")
 			c.Abort()
		}
		return
	}

	srv, err := services.GetOAuthService()

	if err != nil {
		logging.Error("oauth service unavailable", err)
		c.AbortWithStatus(500)
		return
	}

	info, client, err := srv.GetByToken(cookie)

	if err != nil {
		logging.Error("oauth service unavailable", err)
		c.AbortWithStatus(500)
		return
	}

	if info == nil || client == nil {
		c.Redirect(302, "/auth/login")
		c.Abort()
		return
	}

	//does this token have a login scope
	valid := false
	for _, v := range client.ServerScopes {
		if v.ServerId == nil && v.Scope == "login" {
			valid = true
		}
	}

	if !valid {
		c.AbortWithStatus(403)
		return
	}

	err = srv.UpdateExpirationTime(info, 60*time.Minute)
	if err != nil {
		logging.Error("error extending session", err)
		c.AbortWithStatus(500)
		return
	}

	c.Set("client_id", client.ID)
	c.Set("user_id", client.UserID)

	c.Next()
}
