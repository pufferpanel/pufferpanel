/*
 Copyright 2020 Padduck, LLC
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at
  	http://www.apache.org/licenses/LICENSE-2.0
  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package panelmiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3/database"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

var noLogin = []string{"/auth/", "/error/", "/api/config"}
var assetFiles = []string{".js", ".css", ".img", ".ico", ".png", ".gif"}
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
	var user models.User

	cookie, err := c.Cookie("puffer_auth")
	if err == http.ErrNoCookie || cookie == "" {
		// reset err so we don't trip the final error
		// check despite successful auth from header
		err = nil

		//determine if it's an asset, otherwise, we can redirect if it's a GET
		//dev only requirement?
		/*if c.Request.Method == "GET" && strings.Count(c.Request.URL.Path, "/") == 1 {
			for _, v := range assetFiles {
				if strings.HasSuffix(c.Request.URL.Path, v) {
					return
				}
			}
		}*/

		//check for token Auth header
		authHeader := c.Request.Header.Get("Authorization")
		authHeader = strings.TrimSpace(authHeader)

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

		session := parts[1]
		sess, err := ss.Validate(session)

		if err == gorm.ErrRecordNotFound {
			c.Header(WWWAuthenticateHeader, WWWAuthenticateHeaderContents)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		} else if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}

		if sess.UserId != nil {
			user = sess.User
		}
		if sess.ClientId != nil {
			c.Set("client", &sess.Client)
		}
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	} else {
		//pull user from the session
		sess, err := ss.Validate(cookie)

		if err == gorm.ErrRecordNotFound {
			c.Header(WWWAuthenticateHeader, WWWAuthenticateHeaderContents)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		} else if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}

		if sess.UserId != nil {
			user = sess.User
		}
	}

	if response.HandleError(c, err, http.StatusUnauthorized) {
		c.Header(WWWAuthenticateHeader, WWWAuthenticateHeaderContents)
		return
	}

	c.Set("user", &user)
}
