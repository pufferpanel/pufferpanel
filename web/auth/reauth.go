package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/v3/response"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/services"
)

func Reauth(c *gin.Context) {
	res := response.From(c)

	db, err := database.GetConnection()
	if response.HandleError(res, err) {
		return
	}

	user, _ := c.MustGet("user").(*models.User)

	ps := &services.Permission{DB: db}

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
