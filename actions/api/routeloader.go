package api

import (
	"github.com/gobuffalo/buffalo"
	"github.com/pufferpanel/pufferpanel/middleware"
)

func Register(app *buffalo.App) error {
	sub := app.Group("/api")
	sub.Use(middleware.OAuthHandler)
	sub.Use(middleware.RequireAuth(false))
	sub.GET("/server", ServerHandler)
	return nil;
}
