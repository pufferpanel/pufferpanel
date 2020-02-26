package oauth2

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"net/http"
	"strconv"
	"strings"
)

func registerTokens(g *gin.RouterGroup) {
	g.POST("/token", handleTokenRequest)
	g.OPTIONS("/token", response.CreateOptions("POST"))
}

func handleTokenRequest(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-cache")

	if c.GetHeader("Content-Type") != "application/x-www-form-urlencoded" {
		return
	}

	var request oauth2TokenRequest
	err := c.BindWith(&request, binding.FormMultipart)
	if err != nil {
		c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: err.Error()})
		return
	}

	db, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: err.Error()})
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
			_, jwtToken, err := us.Login(user.Email, request.Password)
			if err != nil || jwtToken == "" {
				c.JSON(http.StatusBadRequest, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: "no access"})
				return
			}

			c.Header("Cache-Control", "no-store")
			c.Header("Pragma", "no-cache")
			c.JSON(http.StatusOK, &oauth2TokenResponse{
				AccessToken: jwtToken,
				TokenType:   "Bearer",
				//TODO: Follow OAuth2 more and give better scope information
				Scope: "jwt",
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
