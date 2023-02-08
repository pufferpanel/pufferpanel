package oauth2

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3/middleware/panelmiddleware"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/oauth2"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func registerInfo(g *gin.RouterGroup) {
	g.POST("/introspect", panelmiddleware.NeedsDatabase, handleInfoRequest)
	g.OPTIONS("/introspect", response.CreateOptions("POST"))
}

func handleInfoRequest(c *gin.Context) {
	header := strings.TrimSpace(c.GetHeader("Authorization"))
	if header == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	db := panelmiddleware.GetDatabase(c)
	nodeCreds := parts[1]

	sessionService := &services.Session{DB: db}
	node, err := sessionService.ValidateNode(nodeCreds)
	if err != nil || node.ID == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token := c.DefaultPostForm("token", "")
	if err != nil || token == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	serverId := c.DefaultPostForm("token_type_hint", "")

	infoResponse := &oauth2.TokenInfoResponse{}

	//this token can be one of two types, we need to work out which it is
	session, err := sessionService.Validate(token)
	if err == gorm.ErrRecordNotFound {
		infoResponse.Active = false
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	} else {
		ps := &services.Permission{DB: db}

		var perms []*models.Permissions

		//if this is a client, we only want their perms
		//if there is not a hint, we have to do the kitchen sink
		if serverId == "" {
			if session.ClientId != nil {
				t, err := ps.GetForClient(*session.ClientId)
				if response.HandleError(c, err, http.StatusInternalServerError) {
					return
				}
				perms = append(perms, t...)
			} else if session.UserId != nil {
				t, err := ps.GetForUser(*session.UserId)
				if response.HandleError(c, err, http.StatusInternalServerError) {
					return
				}
				perms = append(perms, t...)
			}
		} else {
			if session.ClientId != nil {
				t, err := ps.GetForClientAndServer(*session.ClientId, &serverId)
				if response.HandleError(c, err, http.StatusInternalServerError) {
					return
				}
				perms = append(perms, t)
			} else if session.UserId != nil {
				t, err := ps.GetForUserAndServer(*session.UserId, &serverId)
				if response.HandleError(c, err, http.StatusInternalServerError) {
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
