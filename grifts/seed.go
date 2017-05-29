package grifts

import (
	"github.com/markbates/grift/grift"
)

var _ = grift.Add("db:seed", func(c *grift.Context) error {
	// Add DB seeding stuff here
	return nil
})
