/*
 Copyright 2019 Padduck, LLC
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

package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/middleware/panelmiddleware"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/oauth2"
	"github.com/pufferpanel/pufferpanel/v3/programs"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"net/http"
	"runtime/debug"
	"strings"
)

func ResponseAndRecover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(error); !ok {
				err = errors.New(pufferpanel.ToString(err))
			}
			response.HandleError(c, err.(error), http.StatusInternalServerError)

			logging.Error.Printf("Error handling route\n%+v\n%s", err, debug.Stack())
			c.Abort()
		}
	}()

	c.Next()
}

func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			logging.Error.Printf("Error handling route\n%+v\n%s", err, debug.Stack())
			c.Abort()
		}
	}()

	c.Next()
}

func RequiresPermission(perm pufferpanel.Scope, needsServer bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		requiresPermission(c, perm, needsServer)
	}
}

func requiresPermission(c *gin.Context, perm pufferpanel.Scope, needsServer bool) {
	//fail-safe in the event something pukes, we don't end up accidently giving rights to something they should not
	actuallyFinished := false
	defer func() {
		if !actuallyFinished && !c.IsAborted() {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()

	//we need to know what type of "instance" we are
	if config.PanelEnabled.Value() {
		panelmiddleware.NeedsDatabase(c)
		if c.IsAborted() {
			return
		}

		//defer to auth middleware to resolve this to a user, and check basic stuff
		panelmiddleware.AuthMiddleware(c)
		if c.IsAborted() {
			return
		}

		//we now have a user and they are allowed to access something, let's confirm they have server access
		serverId := c.Param("serverId")
		if needsServer && serverId == "" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		ginClient, _ := c.Get("client")

		db := panelmiddleware.GetDatabase(c)
		ps := &services.Permission{DB: db}

		var perms []models.Permissions

		//if we're a client, get the client's permissions for this particular resource
		if ginClient != nil {
			client := ginClient.(*models.Client)
			p, err := ps.GetForClientAndServer(client.ID, &serverId)
			if response.HandleError(c, err, http.StatusInternalServerError) {
				return
			}

			perms = append(perms, *p)
			if serverId != "" {
				//if we had a server, also grab global scopes
				p, err = ps.GetForClientAndServer(client.ID, nil)
				if response.HandleError(c, err, http.StatusInternalServerError) {
					return
				}
				perms = append(perms, *p)
			}
		} else {
			user := c.MustGet("user").(*models.User)
			p, err := ps.GetForUserAndServer(user.ID, &serverId)
			if response.HandleError(c, err, http.StatusInternalServerError) {
				return
			}

			perms = append(perms, *p)
			if serverId != "" {
				//if we had a server, also grab global scopes
				p, err = ps.GetForUserAndServer(user.ID, nil)
				if response.HandleError(c, err, http.StatusInternalServerError) {
					return
				}
				perms = append(perms, *p)
			}
		}

		allowed := false
		scopes := make([]pufferpanel.Scope, 0)
		for _, p := range perms {
			for _, v := range p.ToScopes() {
				scopes = append(scopes, v)
				//allow if you have the scope, or if you're a server admin
				if !allowed && (v == perm || v == pufferpanel.ScopeServersAdmin) {
					allowed = true
				}
			}
		}

		if !allowed {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if needsServer && serverId != "" {
			if config.DaemonEnabled.Value() {
				program, err := programs.Get(serverId)
				if response.HandleError(c, err, http.StatusInternalServerError) {
					return
				}

				if program == nil {
					c.AbortWithStatus(http.StatusForbidden)
					return
				}

				c.Set("program", program)
			}

			ss := &services.Server{DB: db}
			server, err := ss.Get(serverId)
			if response.HandleError(c, err, http.StatusInternalServerError) {
				return
			}

			if server == nil {
				//if the server is still null.... how was the client authorized? either way.... 403 it
				c.AbortWithStatus(http.StatusForbidden)
				return
			}

			c.Set("server", server)
		}

		c.Set("scopes", scopes)
		actuallyFinished = true
	} else {
		//we need to ask the panel what permissions this entity has...
		//deferring to the panel's info endpoint to answer most of what we need

		//in this case, server -> program

		//we can short-cut if this is even going to work by seeing if we have a server request
		//if we do, validate it's even one we know of
		serverId := c.Param("id")
		if needsServer && serverId != "" {
			program, err := programs.Get(serverId)
			if response.HandleError(c, err, http.StatusInternalServerError) {
				return
			}

			if program == nil {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}

			c.Set("server", program)
		}

		//request from the auth server the information for this user
		var err error

		info, err := oauth2.GetInfo(GetToken(c), serverId)
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}

		if !info.Active {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		allowed := false
		for _, v := range strings.Split(info.Scope, " ") {
			parts := strings.Split(v, ":")
			if len(parts) != 2 {
				continue
			}
			if parts[0] == serverId && perm.Matches(parts[1]) {
				allowed = true
			}
		}

		if !allowed {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		actuallyFinished = true
	}
}

func GetToken(c *gin.Context) string {
	cookie, err := c.Cookie("puffer_auth")
	if err == http.ErrNoCookie || cookie == "" {
		authHeader := c.Request.Header.Get("Authorization")
		authHeader = strings.TrimSpace(authHeader)

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 {
			return ""
		}

		if parts[0] != "Bearer" || parts[1] == "" {
			return ""
		}

		cookie = parts[1]
	}

	return cookie
}
