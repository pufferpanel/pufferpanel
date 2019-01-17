package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("login", Login)
	rg.POST("login", LoginPost)
	rg.GET("register", Register)
	rg.POST("register", RegisterPost)
}

