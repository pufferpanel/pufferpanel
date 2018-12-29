package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/register", gin.H{})
}