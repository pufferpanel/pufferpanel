package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/login", gin.H{
		"title": "Login",
	})
}

func Logout(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/logout", gin.H{
		"title": "Logout",
	})
}
