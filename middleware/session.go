package middleware

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

type Message struct {
	Message string `json:"msg"`
	Code int `json:"code"`
}

var NOT_AUTHORIZED = Message {
	Message: "not authorized",
	Code: 1,
}

func RequireAuth (shouldRedirect bool) func (next buffalo.Handler) buffalo.Handler {
	return func (next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			id := c.Session().Get("userId")

			_, valid := id.(int)

			if id == nil || !valid {
				if shouldRedirect {
					return c.Redirect(302, "/auth/login")
				} else {
					return c.Render(401, render.JSON(NOT_AUTHORIZED))
				}
			}

			return next(c)
		}
	}
}