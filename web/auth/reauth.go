package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/v3/response"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"github.com/pufferpanel/pufferpanel/v2/web/handlers"
)

func Reauth(c *gin.Context) {
	res := response.From(c)
	db := handlers.GetDatabase(c)
	ps := &services.Permission{DB: db}

	user, _ := c.MustGet("user").(*models.User)

	perms, err := ps.GetForUserAndServer(user.ID, nil)
	if response.HandleError(res, err) {
		return
	}

	session, err := services.GenerateSession(user.ID)
	if response.HandleError(res, err) {
		return
	}

	data := &loginResponse{}
	data.Session = session
	data.Admin = perms.Admin

	res.Data(data)
}
