package oauth2

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pufferpanel/apufferi/v3"
	"github.com/pufferpanel/apufferi/v3/response"
	"github.com/pufferpanel/apufferi/v3/scope"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"strings"
)

func registerTokens(g *gin.RouterGroup) {
	g.POST("/token", handleTokenRequest)
	g.OPTIONS("/token", response.CreateOptions("POST"))
}

func handleTokenRequest(c *gin.Context) {
	response.From(c).Discard()

	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-cache")

	if c.GetHeader("Content-Type") != "application/x-www-form-urlencoded" {
		return
	}

	var request oauth2TokenRequest
	err := c.BindWith(&request, binding.FormMultipart)
	if err != nil {
		c.JSON(400, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: err.Error()})
		return
	}

	switch strings.ToLower(request.GrantType) {
	/*case "client_credentials":
	{

	}*/
	case "password":
		{
			auth := strings.TrimSpace(c.GetHeader("Authorization"))
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				c.Header("WWW-Authenticate", "Bearer")
				c.JSON(401, &oauth2TokenResponse{Error: "invalid_client"})
				return
			}

			//validate this is a bearer token and a good JWT token
			auth = strings.TrimPrefix(auth, "Bearer ")
			token, err := services.ParseToken(auth)
			if err != nil {
				c.JSON(400, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: err.Error()})
				return
			}

			//validate token can auth on behalf of users
			scopes := token.Claims.PanelClaims.Scopes[""]
			if scopes == nil || len(scopes) == 0 || !apufferi.ContainsScope(scopes, scope.OAuth2Auth) {
				c.JSON(400, &oauth2TokenResponse{Error: "unauthorized_client"})
				return
			}

			db, err := database.GetConnection()
			us := &services.User{DB: db}
			ss := &services.Server{DB: db}

			//get user and server information
			parts := strings.SplitN(request.Username, "|", 2)
			user, err := us.GetByEmail(parts[0])
			if err != nil {
				c.JSON(400, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: err.Error()})
				return
			}

			server, err := ss.Get(parts[1])
			if err != nil {
				c.JSON(400, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: err.Error()})
				return
			}

			//confirm user has access to this server
			ps := &services.Permission{DB: db}
			perms, err := ps.GetForUserAndServer(user.ID, &server.Identifier)
			if err != nil {
				c.JSON(400, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: err.Error()})
				return
			}
			if perms.Id == 0 {
				c.JSON(400, &oauth2TokenResponse{Error: "invalid_request", ErrorDescription: "no access"})
				return
			}

			//validate their credentials
			_, jwtToken, err := us.Login(user.Email, request.Password)

			c.Header("Cache-Control", "no-store")
			c.Header("Pragma", "no-cache")
			c.JSON(200, &oauth2TokenResponse{
				AccessToken: jwtToken,
				TokenType:   "Bearer",
				//TODO: Follow OAuth2 more and give better scope information
				Scope: "jwt",
			})
		}
	default:
		c.JSON(400, &oauth2TokenResponse{Error: "unsupported_grant_type"})
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
