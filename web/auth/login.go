package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/v3/response"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"github.com/pufferpanel/pufferpanel/v2/web/handlers"
)

func LoginPost(c *gin.Context) {
	res := response.From(c)
	db := handlers.GetDatabase(c)
	us := &services.User{DB: db}
	ps := &services.Permission{DB: db}

	request := &loginRequest{}

	err := c.BindJSON(request)
	if err != nil {
		res.Message("invalid request").Status(400).Error(err).Fail()
		return
	}

	user, session, err := us.Login(request.Data.Email, request.Data.Password)
	if response.HandleError(res, err) {
		return
	}

	perms, err := ps.GetForUserAndServer(user.ID, nil)
	if response.HandleError(res, err) {
		return
	}

	data := &loginResponse{}
	data.Session = session
	data.Admin = perms.Admin

	res.Data(data)
}

type loginRequest struct {
	Data loginRequestData `json:"data"`
}

type loginRequestData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Session string `json:"session"`
	Admin   bool   `json:"admin,omitempty"`
}
