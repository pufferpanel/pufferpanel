package api

import (
	"github.com/gobuffalo/buffalo"
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
)

func RegisterNodeRoutes (app *buffalo.App) {
	app.PUT("/node", createNode)
	app.PUT("/node/{code}", createNode)
	app.GET("/node/{code}", getNode)
	app.GET("/node", getNodes)
	app.DELETE("/node/{code}", deleteNode)
	app.POST("/node/{code}", editNode)
}

func createNode(c buffalo.Context) (err error) {
	node := models.Node{
	}

	code := c.Param("code")

	err = c.Bind(&node)
	if code != "" {
		node.Code = code
	}

	if SendIfError(c, err) {
		err = nil
		return
	}

	location, err := models.GetLocationById(node.LocationID)

	if SendIfError(c, err) {
		err = nil
		return
	}

	node.Location = &location

	err = node.Save()
	if SendIfError(c, err) {
		err = nil
		return
	} else {
		c.Render(200, render.JSON(node))
	}

	return
}

func getNodes(c buffalo.Context) (err error) {
	nodes, err := models.GetNodes()

	if SendIfError(c, err) {
		err = nil
		return
	} else {
		c.Render(200, render.JSON(nodes))
	}
	return
}

func getNode(c buffalo.Context) (err error) {
	code := c.Param("code")

	node, err := models.GetNodeByCode(code)
	if SendIfError(c, err) {
		err = nil
		return
	}

	if node.ID == uuid.Nil {
		Send404(c, errors.New(ERROR_NONODECODE))
	} else {
		c.Render(200, render.JSON(node))
	}

	return
}

func deleteNode(c buffalo.Context) (err error) {
	code := c.Param("code")

	node, err := models.GetNodeByCode(code)
	if SendIfError(c, err) {
		err = nil
		return
	}
	
	err = node.Delete()

	if SendIfError(c, err) {
		err = nil
	} else if node.ID == uuid.Nil {
		Send404(c, errors.New(ERROR_NONODECODE))
	} else {
		c.Render(200, render.JSON(node))
	}

	return
}

func editNode(c buffalo.Context) (err error) {
	code := c.Param("code")

	newNode := models.Node{}
	err = c.Bind(&newNode)
	if SendIfError(c, err) {
		err = nil
		return
	}

	node, err := models.GetNodeByCode(code)
	if SendIfError(c, err) {
		err = nil
		return
	}

	node.CopyFrom(newNode)
	node.Save()

	return c.Render(200, render.JSON(node))
}

