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

package oauth2

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/middleware"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

func registerTokens(g *gin.RouterGroup) {
	g.POST("/token", middleware.NeedsDatabase, handleTokenRequest)
	g.OPTIONS("/token", response.CreateOptions("POST"))
}

func handleTokenRequest(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-cache")

	var request oauth2TokenRequest
	err := c.MustBindWith(&request, binding.FormPost)
	if err != nil {
		c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: err.Error()})
		return
	}

	db := middleware.GetDatabase(c)
	if db == nil {
		c.JSON(http.StatusInternalServerError, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: "database not available"})
		return
	}

	switch strings.ToLower(request.GrantType) {
	case "client_credentials":
		{
			if strings.HasPrefix(request.ClientId, ".node_") {
				nodeId := strings.TrimPrefix(request.ClientId, ".node_")
				id, err := strconv.Atoi(nodeId)
				if err != nil || id <= 0 {
					if err == nil {
						err = errors.New("node id must be positive")
					}
					c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "invalid_client", ErrorDescription: err.Error()})
					return
				}
				ns := &services.Node{DB: db}
				node, err := ns.Get(uint(id))
				if err == gorm.ErrRecordNotFound {
					c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: "invalid node"})
					return
				} else if err != nil {
					c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: "invalid request"})
					return
				}

				if node.Secret != request.ClientSecret {
					c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "invalid_client"})
					return
				}
				//at this point, we've validate it's a node, we can issue the token
				ts, err := services.GenerateOAuthForNode(node.ID)
				if err != nil {
					c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: err.Error()})
					return
				}

				c.JSON(http.StatusOK, &oauth2TokenResponse{
					AccessToken: ts,
					TokenType:   "Bearer",
					Scope:       string(pufferpanel.ScopeOAuth2Auth),
				})

				return
			} else {
				os := &services.OAuth2{DB: db}
				client, err := os.Get(request.ClientId)
				if err != nil {
					c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: err.Error()})
					return
				}

				if !client.ValidateSecret(request.ClientSecret) {
					c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "invalid_client"})
					return
				}

				//if the client we have has 0 scopes, then we will have to pull the rights the user has
				if len(client.Scopes) == 0 {
					ps := &services.Permission{DB: db}
					var serverId *string
					if client.ServerId != "" {
						serverId = &client.ServerId
					}
					perms, err := ps.GetForUserAndServer(client.UserId, serverId)
					if err != nil {
						c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: err.Error()})
						return
					}
					client.Scopes = perms.ToScopes()
				}

				token, err := services.GenerateOAuthForClient(client)
				if err != nil {
					c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: err.Error()})
					return
				}

				c.JSON(http.StatusOK, &oauth2TokenResponse{
					AccessToken: token,
					TokenType:   "Bearer",
					Scope:       string(pufferpanel.ScopeOAuth2Auth),
				})
				return
			}
		}
	case "password":
		{
			auth := strings.TrimSpace(c.GetHeader("Authorization"))
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				c.Header("WWW-Authenticate", "Bearer")
				c.JSON(http.StatusUnauthorized, &oauth2TokenResponse{Error: "invalid_client"})
				return
			}

			//validate this is a bearer token and a good JWT token
			auth = strings.TrimPrefix(auth, "Bearer ")
			token, err := services.ParseToken(auth)
			if err != nil {
				c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: err.Error()})
				return
			}

			//validate token can auth on behalf of users
			scopes := token.Claims.PanelClaims.Scopes[""]
			if scopes == nil || len(scopes) == 0 || !pufferpanel.ContainsScope(scopes, pufferpanel.ScopeOAuth2Auth) {
				c.JSON(http.StatusOK, &oauth2TokenResponse{Error: "unauthorized_client"})
				return
			}
			us := &services.User{DB: db}
			ss := &services.Server{DB: db}

			//get user and server information
			parts := strings.SplitN(request.Username, "|", 2)
			user, err := us.GetByEmail(parts[0])
			if err != nil {
				c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: err.Error()})
				return
			}

			server, err := ss.Get(parts[1])
			if err != nil {
				c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: err.Error()})
				return
			}

			//confirm user has access to this server
			ps := &services.Permission{DB: db}
			perms, err := ps.GetForUserAndServer(user.ID, &server.Identifier)
			if err != nil {
				c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: err.Error()})
				return
			}
			if perms.ID == 0 || !perms.SFTPServer {
				c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: "no access"})
				return
			}

			//validate their credentials
			user, jwtToken, _, err := us.Login(user.Email, request.Password)
			if err != nil || jwtToken == "" {
				c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: "no access"})
				return
			}

			mappedScopes := make([]string, 0)

			for _, p := range perms.ToScopes() {
				mappedScopes = append(mappedScopes, server.Identifier+":"+string(p))
			}

			c.Header("Cache-Control", "no-store")
			c.Header("Pragma", "no-cache")
			c.JSON(http.StatusOK, &oauth2TokenResponse{
				AccessToken: jwtToken,
				TokenType:   "Bearer",
				Scope:       strings.Join(mappedScopes, " "),
			})
		}
	default:
		c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "unsupported_grant_type"})
	}
}

type oauth2TokenRequest struct {
	GrantType    string `form:"grant_type"`
	ClientId     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
	Username     string `form:"username"`
	Password     string `form:"password"`
}

type oauth2TokenResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	//ExpiresIn        int    `json:"expires_in,omitempty"`
	Scope            string `json:"scope"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}
