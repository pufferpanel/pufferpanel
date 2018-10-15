package main

import (
	"github.com/gin-gonic/gin"
	pufferdHttp "github.com/pufferpanel/apufferi/http"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		pufferdHttp.Respond(c).Send()
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}