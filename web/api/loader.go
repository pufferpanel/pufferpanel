package api

import "github.com/gin-gonic/gin"

func Register(rg *gin.RouterGroup) {
	g := rg.Group("/locations")
	{
		registerLocations(g)
	}

	g = rg.Group("/nodes")
	{
		registerNodes(g)
	}
}