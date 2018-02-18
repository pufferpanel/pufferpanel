package api

import "github.com/gobuffalo/buffalo"

func RegisterUserRoutes (app *buffalo.App) {
	app.POST("/user", createUser)
	app.PUT("/user/{id}", createUserWithId)
	app.GET("/user/{id}", getUser)
	app.GET("/user", getUsers)
	app.DELETE("/user/{id}", deleteUser)
	app.POST("/user/{id}", editUser)
}

func createUser(c buffalo.Context) (err error) {
	return
}

func createUserWithId(c buffalo.Context) (err error) {
	return
}

func getUsers(c buffalo.Context) (err error) {
	return
}

func getUser(c buffalo.Context) (err error) {
	return
}

func deleteUser(c buffalo.Context) (err error) {
	return
}

func editUser(c buffalo.Context) (err error) {
	return
}