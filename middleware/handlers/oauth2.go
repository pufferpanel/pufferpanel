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
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"net/http"
	"strconv"
	"strings"
)

func HasOAuth2Token(c *gin.Context) {
	//if there's a cookie with the token, use that
	cookie, _ := c.Cookie("puffer_auth")
	if cookie != "" && c.Request.Header.Get("Authorization") == "" {
		c.Request.Header.Set("Authorization", "Bearer "+cookie)
	}

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

	token, err := services.ParseToken(parts[1])

	if err != nil || !token.Valid {
		c.Header(WWWAuthenticateHeader, WWWAuthenticateHeaderContents)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("token", token)
	c.Next()
}

func OAuth2Handler(requiredScope pufferpanel.Scope, requireServer bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := database.GetConnection()

		if response.HandleError(c, err, http.StatusInternalServerError) {
		}

		token, ok := c.Get("token")
		if !ok {
			response.HandleError(c, err, http.StatusInternalServerError)
			return
		}
		jwtToken, ok := token.(*pufferpanel.Token)
		if !ok {
			response.HandleError(c, err, http.StatusInternalServerError)
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
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		if requireServer && (server == nil || server.Identifier == "") {
			c.AbortWithStatus(http.StatusForbidden)
			return
		} else if requireServer {
			serverId = server.Identifier
		}

		ui, err := strconv.Atoi(ti.Subject)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user, err := us.GetById(uint(ui))
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		allowed := false

		//if this is an audience of oauth2, we can use token directly
		if ti.Audience == "oauth2" {
			if requiredScope != pufferpanel.ScopeNone {
				scopes := ti.PanelClaims.Scopes[serverId]
				if scopes != nil && pufferpanel.ContainsScope(scopes, requiredScope) {
					allowed = true
				} else {
					//if there isn't a defined rule, is this user an admin?
					scopes := ti.PanelClaims.Scopes[""]
					if scopes != nil && pufferpanel.ContainsScope(scopes, pufferpanel.ScopeServersAdmin) {
						allowed = true
					}
				}
			} else {
				allowed = true
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

			if response.HandleError(c, err, http.StatusInternalServerError) {
				return
			}

			if requiredScope != pufferpanel.ScopeNone {
				if pufferpanel.ContainsScope(perms.ToScopes(), requiredScope) {
					allowed = true
				} else {
					perms, err = ps.GetForUserAndServer(user.ID, nil)
					if response.HandleError(c, err, http.StatusInternalServerError) {
						return
					}
					if pufferpanel.ContainsScope(perms.ToScopes(), pufferpanel.ScopeServersAdmin) {
						allowed = true
					}
				}
			} else {
				allowed = true
			}
		} else {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if !allowed {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Set("server", server)
		c.Set("user", user)
		c.Next()
	}
}
