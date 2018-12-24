package auth

import "github.com/gin-gonic/gin"

func Register(rg *gin.RouterGroup) {
	rg.GET("login", Login)
	rg.GET("logout", Logout)
}

