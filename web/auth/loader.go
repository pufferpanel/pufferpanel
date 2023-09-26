package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.Use(func(c *gin.Context) {
		middleware.ResponseAndRecover(c)
	})
	rg.POST("login", middleware.NeedsDatabase, LoginPost)
	rg.POST("logout", middleware.NeedsDatabase, LogoutPost)
	rg.POST("otp", middleware.NeedsDatabase, OtpPost)
	rg.POST("register", middleware.NeedsDatabase, RegisterPost)
	rg.POST("reauth", middleware.AuthMiddleware, middleware.NeedsDatabase, Reauth)
}
