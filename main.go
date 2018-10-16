package main

import (
	"github.com/gin-gonic/gin"
	pufferdHttp "github.com/pufferpanel/apufferi/http"
	"github.com/pufferpanel/pufferpanel/web"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		pufferdHttp.Respond(c).Send()
	})

	web.RegisterRoutes(r)

	r.Run() // listen and serve on 0.0.0.0:8080
}