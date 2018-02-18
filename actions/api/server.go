package api

import "github.com/gobuffalo/buffalo"

func RegisterServerRoutes (app *buffalo.App) {
	app.POST("/server", createServer)
	app.PUT("/server/{id}", createServerWithId)
	app.GET("/server/{id}", getServer)
	app.GET("/server", getServers)
	app.DELETE("/server/{id}", deleteServer)
	app.POST("/server/{id}", editServer)
}

func createServer(c buffalo.Context) (err error) {
	return
}

func createServerWithId(c buffalo.Context) (err error) {
	return
}

func getServers(c buffalo.Context) (err error) {
	return
}

func getServer(c buffalo.Context) (err error) {
	return
}

func deleteServer(c buffalo.Context) (err error) {
	return
}

func editServer(c buffalo.Context) (err error) {
	return
}
