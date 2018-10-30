package oauth2

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/oauth2"
	"github.com/pufferpanel/pufferpanel/shared"
	"gopkg.in/oauth2.v3/manage"
)

func registerTokens(g *gin.RouterGroup) {
	handle()
	g.POST("/request", shared.NotImplemented)
	g.OPTIONS("/request", shared.CreateOptions("POST"))

	g.POST("/info", shared.NotImplemented)
	g.OPTIONS("/info", shared.CreateOptions("POST"))
}

func handle() {
	manager := manage.NewDefaultManager()
	manager.MapClientStorage(&oauth2.ClientStore{})
}