package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.Use(func(c *gin.Context) {
		middleware.ResponseAndRecover(c)
	})
	rg.POST("login", LoginPost)
	rg.POST("register", RegisterPost)
}
