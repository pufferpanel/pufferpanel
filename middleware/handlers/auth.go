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

package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"net/http"
	"strconv"
	"strings"
)

var noLogin = []string{"/auth/", "/error/", "/daemon/", "/api/"}
var assetFiles = []string{".js", ".css", ".img", ".ico", ".png", ".gif"}
var overrideRequireLogin = []string{"/auth/reauth"}

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

		c.Header(WWWAuthenticateHeader, WWWAuthenticateHeaderContents)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	db, err := database.GetConnection()
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	token, err := services.ParseToken(cookie)

	if response.HandleError(c, err, http.StatusUnauthorized) {
		c.Header(WWWAuthenticateHeader, WWWAuthenticateHeaderContents)
		return
	}
	if !token.Valid {
		c.Header(WWWAuthenticateHeader, WWWAuthenticateHeaderContents)
		response.HandleError(c, pufferpanel.ErrTokenInvalid, http.StatusUnauthorized)
		return
	}

	us := services.User{DB: db}
	id, err := strconv.Atoi(token.Claims.Subject)
	if response.HandleError(c, err, http.StatusUnauthorized) {
		c.Header(WWWAuthenticateHeader, WWWAuthenticateHeaderContents)
		return
	}
	user, err := us.GetById(uint(id))
	if response.HandleError(c, err, http.StatusUnauthorized) {
		c.Header(WWWAuthenticateHeader, WWWAuthenticateHeaderContents)
		return
	}

	c.Set("user", user)
	c.Next()
}
