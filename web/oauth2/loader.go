package oauth2

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/oauth2"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/spf13/cast"
	"net/http"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/token", setHeaders, recovery, middleware.NeedsDatabase, handleTokenRequest)
	rg.OPTIONS("/token", response.CreateOptions("POST"))
	rg.GET("openid-configuration", TokenServiceGetPublicKey)
}

func setHeaders(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-cache")
}

func recovery(c *gin.Context) {
	//override the recovery route, as we need to change the type returned
	defer func() {
		if err := recover(); err != nil {
			var msg string
			if e, ok := err.(error); ok {
				msg = e.Error()
			} else if e, ok := cast.ToStringE(msg); ok == nil {
				msg = e
			} else {
				msg = fmt.Sprintf("%v", err)
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, &oauth2.ErrorResponse{
				Error:            "internal_error",
				ErrorDescription: msg,
			})
		}
	}()
	c.Next()
}
