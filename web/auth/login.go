package auth

import (
	"github.com/gin-gonic/gin"
	builder "github.com/pufferpanel/apufferi/http"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/shared"
)

func LoginPost(c *gin.Context) {
	response := builder.Respond(c)
	defer response.Send()

	request := &loginRequest{}

	err := c.BindJSON(request)
	if err != nil {
		response.Message("invalid request").Status(400).Data(err.Error()).Fail()
		return
	}

	us, err := services.GetUserService()
	if shared.HandleError(response, err) {
		return
	}

	session, err := us.Login(request.Data.Email, request.Data.Password)

	if err != nil {
		response.Message(err.Error()).Status(400).Fail()
	} else {
		response.Data(session)
	}
}

type loginRequest struct {
	Data loginRequestData `json:"data"`
}

type loginRequestData struct {
	Email string	`json:"email"`
	Password string	`json:"password"`
}