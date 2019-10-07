package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v2/web/handlers"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.Use(func(c *gin.Context) {
		middleware.ResponseAndRecover(c)
	})
	rg.POST("login", LoginPost)
	rg.POST("register", RegisterPost)
	rg.POST("reauth", handlers.AuthMiddleware, Reauth)
}
