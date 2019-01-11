package server

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/web/handlers"
)

func RegisterRoutes(e *gin.RouterGroup) {
	e.GET("", handlers.AuthMiddleware, Index)
}
