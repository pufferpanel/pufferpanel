package server

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(e *gin.RouterGroup) {
	e.GET("", /*handlers.AuthMiddleware,*/ Index)
}
