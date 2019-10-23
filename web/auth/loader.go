package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/v4/middleware"
	"github.com/pufferpanel/pufferpanel/v2/web/handlers"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.Use(func(c *gin.Context) {
		middleware.ResponseAndRecover(c)
	})
	rg.POST("login", handlers.NeedsDatabase, LoginPost)
	rg.POST("register", handlers.NeedsDatabase, RegisterPost)
	rg.POST("reauth", handlers.AuthMiddleware, handlers.NeedsDatabase, Reauth)
}
