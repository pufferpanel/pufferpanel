package api

import (
	"github.com/gobuffalo/buffalo"
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/uuid"
	"errors"
)

func RegisterLocationRoutes (app *buffalo.App) {
	app.PUT("/location/{code}", createLocation)
	app.GET("/location/{code}", getLocation)
	app.GET("/location", getLocations)
	app.DELETE("/location/{code}", deleteLocation)
	app.POST("/location/{code}", editLocationName)
}

func createLocation(c buffalo.Context) (err error) {
	code := c.Param("code")

	location := models.Location{}

	err = c.Bind(&location)

	location.Code = code

	err = location.Save()

	if SendIfError(c, err) {
		err = nil
	} else {
		c.Render(200, render.JSON(location))
	}
	return
}

func getLocations(c buffalo.Context) (err error) {
	locations, err := models.GetLocations()

	if SendIfError(c, err) {
		err = nil
	} else {
		c.Render(200, render.JSON(locations))
	}
	return
}

func getLocation(c buffalo.Context) (err error) {
	code := c.Param("code")

	location, err := models.GetLocationByCode(code)
	if SendIfError(c, err) {
		err = nil
	} else if location.ID == uuid.Nil {
		Send404(c, errors.New(ERROR_NOLOCATIONCODE))
		err = nil
	} else {
		c.Render(200, render.JSON(location))
	}

	return
}

func deleteLocation(c buffalo.Context) (err error) {
	code := c.Param("code")

	location, err := models.GetLocationByCode(code)

	if SendIfError(c, err) {
		err = nil
		return
	}

	err = location.Delete()

	if location.ID == uuid.Nil {
		Send404(c, errors.New(ERROR_NOLOCATIONCODE))
		err = nil
	} else {
		c.Render(200, render.JSON(location))
	}
	return
}

func editLocationName(c buffalo.Context) (err error) {
	code := c.Param("code")

	location, err := models.GetLocationByCode(code)

	if location.ID == uuid.Nil {
		Send404(c, errors.New(ERROR_NOLOCATIONCODE))
		err = nil
		return
	}

	name := models.Location{}

	err = c.Bind(&name)

	if SendIfError(c, err) {
		err = nil
		return
	}

	location.Name = name.Name
	err = location.Save()

	if location.ID == uuid.Nil {
		Send404(c, errors.New(ERROR_NOLOCATIONCODE))
		err = nil
	} else {
		c.Render(200, render.JSON(location))
	}

	return
}