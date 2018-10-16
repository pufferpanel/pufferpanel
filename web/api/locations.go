package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/web/shared"
)

func registerLocations(g *gin.RouterGroup) {
	g.Handle("GET", "", shared.NotImplemented)
	g.Handle("OPTIONS", "", shared.CreateOptions("GET"))

	g.Handle("PUT", "/:id", shared.NotImplemented)
	g.Handle("GET", "/:id", shared.NotImplemented)
	g.Handle("POST", "/:id", shared.NotImplemented)
	g.Handle("DELETE", "/:id", shared.NotImplemented)
	g.Handle("OPTIONS", "/:id", shared.CreateOptions("PUT", "GET", "POST", "DELETE"))
}