package actions

import "github.com/gobuffalo/buffalo"

func Register(app *buffalo.App) error {
	app.GET("/", HomeHandler)
	return nil;
}