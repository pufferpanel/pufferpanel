package api

import "github.com/gobuffalo/buffalo"

func Register(app *buffalo.App) error {
	sub := app.Group("/server")
	sub.GET("/", ServerHandler)
	return nil;
}
