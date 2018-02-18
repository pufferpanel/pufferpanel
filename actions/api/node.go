package api

import "github.com/gobuffalo/buffalo"

func RegisterNodeRoutes (app *buffalo.App) {
	app.POST("/node", createNode)
	app.PUT("/node/{id}", createNodeWithId)
	app.GET("/node/{id}", getNode)
	app.GET("/node", getNodes)
	app.DELETE("/node/{id}", deleteNode)
	app.POST("/node/{id}", editNode)
}

func createNode(c buffalo.Context) (err error) {
	return
}

func createNodeWithId(c buffalo.Context) (err error) {
	return
}

func getNodes(c buffalo.Context) (err error) {
	return
}

func getNode(c buffalo.Context) (err error) {
	return
}

func deleteNode(c buffalo.Context) (err error) {
	return
}

func editNode(c buffalo.Context) (err error) {
	return
}

