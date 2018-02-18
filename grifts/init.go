package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/pufferpanel/pufferpanel/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
