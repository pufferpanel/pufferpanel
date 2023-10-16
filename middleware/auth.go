package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3/database"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

var noLogin = []string{"/auth/", "/error/", "/api/config"}
var overrideRequireLogin = []string{"/auth/reauth", "/auth/logout"}

const WWWAuthenticateHeader = "WWW-Authenticate"
const WWWAuthenticateHeaderContents = "Bearer realm=\"\""

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
				return
			}
		}
	}

	db, err := database.GetConnection()
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	ss := services.Session{DB: db}

	var token string

	//order of priority, use auth headers first
	//check for token Auth header
	authHeader := c.Request.Header.Get("Authorization")
	authHeader = strings.TrimSpace(authHeader)

	if authHeader == "" {
		token, err = c.Cookie("puffer_auth")

		if errors.Is(err, http.ErrNoCookie) || token == "" {
			c.Header(WWWAuthenticateHeader, WWWAuthenticateHeaderContents)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	} else {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 {
			c.Header(WWWAuthenticateHeader, WWWAuthenticateHeaderContents)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if parts[0] != "Bearer" || parts[1] == "" {
			c.Header(WWWAuthenticateHeader, WWWAuthenticateHeaderContents)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token = parts[1]
	}

	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	//pull user from the session
	sess, err := ss.Validate(token)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.Header(WWWAuthenticateHeader, WWWAuthenticateHeaderContents)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	if sess.UserId != nil {
		c.Set("user", &sess.User)
	}
	if sess.ClientId != nil {
		c.Set("client", &sess.Client)
	}
}
