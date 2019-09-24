package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/v3/response"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/oauth2/claims"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"strings"
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
	if response.HandleError(res, err) {
		c.Abort()
	}

	token, err := services.ParseToken(cookie, &claims.UserClaims{})

	if err != nil || !token.Valid {
		c.Redirect(302, "/auth/login")
		c.Abort()
		return
	}

	userClaims := token.Claims.(*claims.UserClaims)

	srv := services.GetOAuth(db)
	client, _, err := srv.GetByUser(&models.User{ID: userClaims.UserId})

	if response.HandleError(res, err) {
		c.Redirect(302, "/auth/login")
		c.Abort()
		return
	}

	if client == nil {
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

	c.Set("client_id", client.ID)
	c.Set("user_id", client.UserID)
	c.Set("user", &client.User)

	c.Next()
}
