package api

import (
	"github.com/gobuffalo/buffalo"
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/gobuffalo/buffalo/render"
)

//Lists all servers the given context has read access to (either owner or subuser)
func ServerHandler(g buffalo.Context) error {
	user := &models.User{
		ID: g.Session().Get("userId").(int),
	}
	servers := &models.Servers{}
	err := models.DB.BelongsTo(user).All(servers)
	if err != nil {
		return err
	}
	return g.Render(200, render.JSON(servers))
}
