package api

import "github.com/gobuffalo/buffalo"

func Register(app *buffalo.App) error {
	sub := app.Group("/api")
	sub.GET("/server", ServerHandler)
	return nil;
}
