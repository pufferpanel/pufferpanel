package api

import (
	"github.com/gobuffalo/buffalo"
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/gobuffalo/buffalo/render"
)

type newLocation struct {
	Name string
}

type editLocation struct {
	Name string
}

func RegisterLocationRoutes (app *buffalo.App) {
	app.PUT("/location/{code}", createLocation)
	app.GET("/location/{code}", getLocation)
	app.GET("/location", getLocations)
	app.DELETE("/location/{code}", deleteLocation)
	app.PUT("/location/{code}/name", editLocationName)
}

func createLocation(c buffalo.Context) (err error) {
	code := c.Param("code")

	loc := newLocation{}

	err = c.Bind(&loc)

	location, err := models.CreateLocation(code, loc.Name)

	if err != nil {
		return
	}

	c.Render(200, render.JSON(location))
	return
}

func getLocations(c buffalo.Context) (err error) {
	locations, err := models.GetLocations()

	if err != nil {
		return
	}

	c.Render(200, render.JSON(locations))
	return
}

func getLocation(c buffalo.Context) (err error) {
	code := c.Param("code")

	location, err := models.GetLocationByCode(code)
	if err != nil {
		return

	}
	c.Render(200, render.JSON(location))
	return
}

func deleteLocation(c buffalo.Context) (err error) {
	code := c.Param("code")

	location, err := models.GetLocationByCode(code)

	if err != nil {
		return
	}

	err = location.Delete()
	return
}

func editLocationName(c buffalo.Context) (err error) {
	code := c.Param("code")

	location, err := models.GetLocationByCode(code)

	if err != nil {
		return
	}

	name := editLocation{}

	err = c.Bind(&name)

	if err != nil {
		return
	}

	location.Name = name.Name
	location.Save()

	return
}