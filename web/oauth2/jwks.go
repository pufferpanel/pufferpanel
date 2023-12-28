package oauth2

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"net/http"
)

func TokenServiceGetPublicKey(c *gin.Context) {
	ts, err := services.NewTokenService()
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}
	rawJWKS, err := ts.GetTokenStore().JSONPublic(context.Background())
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}
	c.JSON(http.StatusOK, rawJWKS)
}
