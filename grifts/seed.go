package grifts

import (
	"github.com/markbates/grift/grift"
)

var _ = grift.Add("db:seed", func(c *grift.Context) error {
	return nil
})
