package web

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/web/api"
)

func RegisterRoutes(e *gin.Engine) {
	apiGroup := e.Group("/api")
	{
		api.Register(apiGroup)
	}
}