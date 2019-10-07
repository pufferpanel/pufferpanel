package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/v3/response"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"strings"
)

var noLogin = []string{"/auth/", "/error/", "/daemon/", "/api/"}
var assetFiles = []string{".js", ".css", ".img", ".ico", ".png", ".gif"}
var overrideRequireLogin = []string{"/auth/reauth"}

func AuthMiddleware(c *gin.Context) {
	for _, v := range noLogin {
		if strings.HasPrefix(c.Request.URL.Path, v) {
			//and now we see if it's actually one we override
			skip := false
			for _, o := range overrideRequireLogin {
				if o == c.Request.URL.Path {
					skip = true
					break
				}
			}
			if !skip {
				c.Next()
				return
			}
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

	token, err := services.ParseToken(cookie)

	if err != nil || !token.Valid {
		c.Redirect(302, "/auth/login")
		c.Abort()
		return
	}

	if response.HandleError(res, err) {
		c.Redirect(302, "/auth/login")
		c.Abort()
		return
	}

	us := services.User{DB: db}
	user, err := us.Get(token.Claims.Subject)
	if response.HandleError(res, err) {
		c.Redirect(302, "/auth/login")
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Next()
}
