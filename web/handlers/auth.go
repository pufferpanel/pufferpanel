package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferpanel/services"
	"time"
)

func AuthMiddleware(c *gin.Context) {

	cookie, err := c.Cookie("puffer_auth")

	if err != nil || cookie == "" {
		c.SetCookie("puffer_auth", "", 3600, "/", "", false, false)
		c.Redirect(302, "/auth/login")
		c.Abort()
		return
	}

	srv, err := services.GetOAuthService()

	if err != nil {
		logging.Error("oauth service unavailable", err)
		c.Redirect(302, "/error/500")
		c.Abort()
		return
	}

	info, client, err := srv.GetByToken(cookie)

	if err != nil {
		logging.Error("oauth service unavailable", err)
		c.Redirect(302, "/error/500")
		c.Abort()
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
		c.Redirect(302, "/auth/login")
		c.Abort()
		return
	}

	err = srv.UpdateExpirationTime(info, 60*time.Minute)
	if err != nil {
		logging.Error("error extending session", err)
		c.Redirect(500, "/error/500")
		c.Abort()
		return
	}

	c.Set("client_id", client.ID)
	c.Set("user_id", client.UserID)

	c.Next()
}
