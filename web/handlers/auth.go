package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/response"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/shared"
	"strings"
	"time"
)

var noLogin = []string{"/auth/", "/error/", "/daemon/", "/api/"}
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
		if c.Request.Method == "GET" && strings.Count(c.Request.URL.Path, "/") == 1 {
			for _, v := range assetFiles {
				if strings.HasSuffix(c.Request.URL.Path, v) {
					return
				}
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

	res := response.From(c)

	db, err := database.GetConnection()
	if shared.HandleError(res, err) {
		c.Abort()
	}

	srv := services.GetOAuth(db)

	info, client, err := srv.GetByToken(cookie)

	if shared.HandleError(res, err) {
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
		c.AbortWithStatus(403)
		return
	}

	err = srv.UpdateExpirationTime(info, 60*time.Minute)
	if shared.HandleError(res, err) {
		c.Abort()
		return
	}

	c.Set("client_id", client.ID)
	c.Set("user_id", client.UserID)
	c.Set("user", &client.User)

	c.Next()
}
