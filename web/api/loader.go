package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
)

const MaxPageSize = 100
const DefaultPageSize = 20

func RegisterRoutes(rg *gin.RouterGroup) {

	rg.Use(func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Next()
	})

	rg.Use(middleware.ResponseAndRecover)
	rg.Use(middleware.NeedsDatabase)
	rg.Use(middleware.AuthMiddleware)
	rg.Use(middleware.AddVersionHeader)
	registerNodes(rg.Group("/nodes"))
	registerServers(rg.Group("/servers"))
	registerUsers(rg.Group("/users"))
	registerTemplates(rg.Group("/templates"))
	registerSelf(rg.Group("/self"))
	registerSettings(rg.Group("/settings"))
	registerUserSettings(rg.Group("/userSettings"))

	rg.GET("/config", panelConfig)
}
