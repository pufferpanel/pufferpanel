package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("login", LoginPost)
	rg.POST("register", RegisterPost)
}

