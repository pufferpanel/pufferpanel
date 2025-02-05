package oauth2

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/oauth2"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/scopes"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"net/http"
	"strings"
	"time"
)

const expiresIn = int64(time.Hour / time.Second)

// @Summary Authenticate
// @Description Get a OAuth2 token to consume this API
// @Param request formData OAuth2TokenRequest true "OAuth2 token request"
// @Success 200 {object} oauth2.TokenResponse
// @Failure 400 {object} oauth2.ErrorResponse
// @Failure 401 {object} oauth2.ErrorResponse
// @Failure 500 {object} oauth2.ErrorResponse
// @Accept x-www-form-urlencoded
// @Produce json
// @Router /oauth2/token [post]
func handleTokenRequest(c *gin.Context) {
	var request OAuth2TokenRequest
	err := c.MustBindWith(&request, binding.FormPost)
	if err != nil {
		c.JSON(http.StatusBadRequest, &oauth2.ErrorResponse{Error: "invalid_request", ErrorDescription: err.Error()})
		return
	}

	db := middleware.GetDatabase(c)
	if db == nil {
		c.JSON(http.StatusInternalServerError, &oauth2.ErrorResponse{Error: "invalid_request", ErrorDescription: "database not available"})
		return
	}

	session := &services.Session{DB: db}

	switch strings.ToLower(request.GrantType) {
	case "client_credentials":
		{
			os := &services.OAuth2{DB: db}
			client, err := os.Get(request.ClientId)
			if err != nil {
				c.JSON(http.StatusBadRequest, &oauth2.ErrorResponse{Error: "invalid_request", ErrorDescription: err.Error()})
				return
			}

			if !client.ValidateSecret(request.ClientSecret) {
				c.JSON(http.StatusBadRequest, &oauth2.ErrorResponse{Error: "invalid_client"})
				return
			}

			token, err := session.CreateForClient(client)
			if err != nil {
				if response.HandleError(c, err, http.StatusInternalServerError) {
					return
				}
			}

			var serverId string
			if client.Server != nil {
				serverId = client.Server.Identifier
			}

			var allScopes []string
			ps := &services.Permission{DB: db}
			perms, err := ps.GetForUserAndServer(client.UserId, serverId)
			if response.HandleError(c, err, http.StatusInternalServerError) {
				return
			}

			if !scopes.ContainsScope(perms.Scopes, scopes.ScopeLogin) {
				//because servers don't have an explicit login scope, we need to check the root user
				if serverId == "" {
					c.AbortWithStatus(http.StatusForbidden)
					return
				}

				userPerms, err := ps.GetForUserAndServer(client.UserId, "")
				if response.HandleError(c, err, http.StatusInternalServerError) {
					return
				}

				if !scopes.ContainsScope(userPerms.Scopes, scopes.ScopeLogin) {
					c.AbortWithStatus(http.StatusForbidden)
					return
				}
			}
			for _, v := range perms.Scopes {
				allScopes = append(allScopes, v.String())
			}

			c.JSON(http.StatusOK, &oauth2.TokenResponse{
				AccessToken: token,
				TokenType:   "Bearer",
				Scope:       strings.Join(allScopes, " "),
				ExpiresIn:   expiresIn,
			})
			return
		}
	case "password":
		{
			auth := strings.TrimSpace(c.GetHeader("Authorization"))
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				c.Header("WWW-Authenticate", "Bearer")
				c.JSON(http.StatusUnauthorized, &oauth2.ErrorResponse{Error: "invalid_client"})
				return
			}

			//validate this is a bearer token and a good JWT token
			auth = strings.TrimPrefix(auth, "Bearer ")
			node, err := session.ValidateNode(auth)
			if err != nil {
				c.JSON(http.StatusBadRequest, &oauth2.ErrorResponse{Error: "invalid_request", ErrorDescription: err.Error()})
				return
			}

			us := &services.User{DB: db}
			ss := &services.Server{DB: db}

			//get user and server information
			parts := strings.SplitN(request.Username, "#", 2)
			if len(parts) != 2 {
				c.JSON(http.StatusBadRequest, &oauth2.ErrorResponse{Error: "invalid_request", ErrorDescription: "bad username"})
				return
			}
			user, err := us.GetByEmail(parts[0])
			if err != nil {
				c.JSON(http.StatusBadRequest, &oauth2.ErrorResponse{Error: "invalid_request", ErrorDescription: err.Error()})
				return
			}

			server, err := ss.Get(parts[1])
			if err != nil {
				c.JSON(http.StatusBadRequest, &oauth2.ErrorResponse{Error: "invalid_request", ErrorDescription: err.Error()})
				return
			}

			//ensure the node asking for the credential check is where this server is
			if server.Node.ID != node.ID {
				c.JSON(http.StatusBadRequest, &oauth2.ErrorResponse{Error: "invalid_request", ErrorDescription: "no access"})
				return
			}

			//validate their credentials
			var token string
			user, _, err = us.ValidateLogin(user.Email, request.Password)
			if err != nil {
				c.JSON(http.StatusBadRequest, &oauth2.ErrorResponse{Error: "invalid_request", ErrorDescription: "no access"})
				return
			}

			//confirm user has access to this server
			ps := &services.Permission{DB: db}
			allowed, err := ps.HasPermission(user.ID, server.Identifier, scopes.ScopeServerSftp)
			if err != nil {
				c.JSON(http.StatusBadRequest, &oauth2.ErrorResponse{Error: "invalid_request", ErrorDescription: err.Error()})
				return
			}

			if !allowed {
				c.JSON(http.StatusBadRequest, &oauth2.ErrorResponse{Error: "invalid_request", ErrorDescription: "no access"})
				return
			}

			//at this point, their login credentials were valid, and we need to shortcut because otp
			sessionService := &services.Session{DB: db}
			token, err = sessionService.CreateForUser(user)
			if err != nil {
				logging.Error.Printf("Error generating token: %s", err.Error())
				c.JSON(http.StatusBadRequest, &oauth2.ErrorResponse{Error: "invalid_request", ErrorDescription: "no access"})
				return
			}

			c.JSON(http.StatusOK, &oauth2.TokenResponse{
				AccessToken: token,
				TokenType:   "Bearer",
				Scope:       server.Identifier + ":" + scopes.ScopeServerSftp.String(),
				ExpiresIn:   expiresIn,
			})
		}
	default:
		c.JSON(http.StatusBadRequest, &oauth2.ErrorResponse{Error: "unsupported_grant_type"})
	}
}

type OAuth2TokenRequest struct {
	GrantType    string `form:"grant_type"`
	ClientId     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
	Username     string `form:"username"`
	Password     string `form:"password"`
} //@name OAuth2TokenRequest
