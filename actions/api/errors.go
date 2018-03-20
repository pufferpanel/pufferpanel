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
		sendCodeWithError(c, 500, err)
		return true
	}
	return false
}

func Send404(c buffalo.Context, err error) {
	sendCodeWithError(c, 404, err)
}

func sendCodeWithError(c buffalo.Context, code int, err error) {
	errRes := make(map[string]string)
	errRes["err"] = err.Error()
	c.Render(code, render.JSON(errRes))
}