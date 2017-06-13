package middleware

import (
	"github.com/gobuffalo/buffalo"
)

func OAuthHandler(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		//todo: add oauth as method of authentication

		return next(c)
	}
}