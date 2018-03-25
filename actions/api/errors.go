package api

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

const (
	ERROR_NONODECODE = "no node with given code"
	ERROR_NOLOCATIONCODE = "no location with given code"
)

func SendIfError(c buffalo.Context, err error) bool {
	if err != nil {
		SendCodeWithError(c, 500, err)
		return true
	}
	return false
}

func Send404(c buffalo.Context, err error) {
	SendCodeWithError(c, 404, err)
}

func SendCodeWithError(c buffalo.Context, code int, err error) {
	errRes := make(map[string]string)
	errRes["err"] = err.Error()
	c.Render(code, render.JSON(errRes))
}