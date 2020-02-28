package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v2/middleware"
	"github.com/pufferpanel/pufferpanel/v2/middleware/handlers"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.Use(func(c *gin.Context) {
		middleware.ResponseAndRecover(c)
	})
	rg.POST("login", middleware.NeedsDatabase, LoginPost)
	rg.POST("register", middleware.NeedsDatabase, RegisterPost)
	rg.POST("reauth", handlers.AuthMiddleware, middleware.NeedsDatabase, Reauth)
}
