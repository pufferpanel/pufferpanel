package api

import (
	"github.com/gobuffalo/buffalo"
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/uuid"
)

type newLocation struct {
	Name string
}

type editLocation struct {
	Name string
}

const (
	ERROR_NOLOCATIONID = "no location with given id"
	ERROR_NOLOCATIONCODE = "no location with given code"
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

	loc := newLocation{}

	err = c.Bind(&loc)

	location, err := models.CreateLocation(code, loc.Name)

	if err == nil {
		err = location.Save()
	}

	if err != nil {
		errRes := make(map[string]string)
		errRes["err"] = err.Error()
		c.Render(500, render.JSON(errRes))
		err = nil
		return
	}

	c.Render(200, render.JSON(location))
	return
}

func getLocations(c buffalo.Context) (err error) {
	locations, err := models.GetLocations()

	if err != nil {
		errRes := make(map[string]string)
		errRes["err"] = err.Error()
		c.Render(500, render.JSON(errRes))
		err = nil
		return
	}

	c.Render(200, render.JSON(locations))
	return
}

func getLocation(c buffalo.Context) (err error) {
	code := c.Param("code")

	location, err := models.GetLocationByCode(code)
	if err != nil {
		errRes := make(map[string]string)
		errRes["err"] = err.Error()
		c.Render(500, render.JSON(errRes))
		err = nil
		return
	}

	if location.ID == uuid.Nil {
		err := make(map[string]string)
		err["err"] = ERROR_NOLOCATIONCODE
		c.Render(404, render.JSON(err))
		err = nil
	} else {
		c.Render(200, render.JSON(location))
	}

	return
}

func deleteLocation(c buffalo.Context) (err error) {
	code := c.Param("code")

	location, err := models.GetLocationByCode(code)

	if err != nil {
		errRes := make(map[string]string)
		errRes["err"] = err.Error()
		c.Render(500, render.JSON(errRes))
		err = nil
		return
	}

	err = location.Delete()

	if location.ID == uuid.Nil {
		err := make(map[string]string)
		err["err"] = ERROR_NOLOCATIONCODE
		c.Render(404, render.JSON(err))
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
		errRes := make(map[string]string)
		errRes["err"] = ERROR_NOLOCATIONCODE
		c.Render(404, render.JSON(errRes))
		err = nil
		return
	}

	name := editLocation{}

	err = c.Bind(&name)

	if err != nil {
		errRes := make(map[string]string)
		errRes["err"] = err.Error()
		c.Render(500, render.JSON(errRes))
		err = nil
		return
	}

	location.Name = name.Name
	err = location.Save()

	if err != nil {
		errRes := make(map[string]string)
		errRes["err"] = err.Error()
		c.Render(500, render.JSON(errRes))
		err = nil
		return
	} else {
		c.Render(200, render.JSON(location))
	}

	return
}