package oauth2

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/oauth2"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func registerInfo(g *gin.RouterGroup) {
	g.POST("/introspect", middleware.NeedsDatabase, handleInfoRequest)
	g.OPTIONS("/introspect", response.CreateOptions("POST"))
}

// @Summary Get info
// @Description Get information about an OAuth2 token
// @Accept x-www-form-urlencoded
// @Param token formData string true "OAuth2 token"
// @Param token_type_hint formData string false "Hint for how the token might be used"
// @Success 200 {object} oauth2.TokenInfoResponse
// @Failure 400 {object} oauth2.ErrorResponse
// @Failure 401 {object} oauth2.ErrorResponse
// @Failure 500 {object} oauth2.ErrorResponse
// @Router /oauth2/introspect [post]
// @Security OAuth2Application[none]
func handleInfoRequest(c *gin.Context) {
	header := strings.TrimSpace(c.GetHeader("Authorization"))
	if header == "" {
		c.Header("WWW-Authenticate", "Bearer")
		c.AbortWithStatusJSON(http.StatusUnauthorized, &oauth2.ErrorResponse{Error: "invalid_request"})
		return
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.Header("WWW-Authenticate", "Bearer")
		c.AbortWithStatusJSON(http.StatusUnauthorized, &oauth2.ErrorResponse{Error: "invalid_request"})
		return
	}

	db := middleware.GetDatabase(c)
	nodeAuthToken := parts[1]

	sessionService := &services.Session{DB: db}
	node, err := sessionService.ValidateNode(nodeAuthToken)
	if err != nil || node.ID == 0 {
		c.Header("WWW-Authenticate", "Bearer")
		c.AbortWithStatusJSON(http.StatusUnauthorized, &oauth2.ErrorResponse{Error: "invalid_request"})
		return
	}

	token := c.DefaultPostForm("token", "")
	if err != nil || token == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, &oauth2.ErrorResponse{Error: "invalid_request"})
		return
	}

	serverId := c.DefaultPostForm("token_type_hint", "")

	infoResponse := &oauth2.TokenInfoResponse{}

	//this token can be one of two types, we need to work out which it is
	session, err := sessionService.Validate(token)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		infoResponse.Active = false
	} else if err != nil {
		panic(err) //trigger standard recover
		return
	} else {
		ps := &services.Permission{DB: db}

		var perms []*models.Permissions

		//if this is a client, we only want their perms
		//if there is not a hint, we have to do the kitchen sink
		if serverId == "" {
			if session.ClientId != nil {
				t, err := ps.GetForClient(*session.ClientId)
				if err != nil {
					panic(err) //trigger standard recover
					return
				}
				perms = append(perms, t...)
			} else if session.UserId != nil {
				t, err := ps.GetForUser(*session.UserId)
				if err != nil {
					panic(err) //trigger standard recover
					return
				}
				perms = append(perms, t...)
			}
		} else {
			if session.ClientId != nil {
				t, err := ps.GetForClientAndServer(*session.ClientId, &serverId)
				if err != nil {
					panic(err) //trigger standard recover
					return
				}
				perms = append(perms, t)
			} else if session.UserId != nil {
				t, err := ps.GetForUserAndServer(*session.UserId, &serverId)
				if err != nil {
					panic(err) //trigger standard recover
					return
				}
				perms = append(perms, t)
			}
		}

		var scopes []string

		for _, v := range perms {
			for _, d := range v.ToScopes() {
				//limit scopes we return to just the node that this server tracks
				//this helps avoid giant lists for this node
				if v.ServerIdentifier != nil && v.Server.NodeID == node.ID {
					scopes = append(scopes, fmt.Sprintf("%s:%s", *v.ServerIdentifier, d))
				} else {
					scopes = append(scopes, fmt.Sprintf(":%s", d))
				}
			}
		}

		infoResponse.Scope = strings.Join(scopes, " ")
	}

	c.JSON(http.StatusOK, infoResponse)
}
