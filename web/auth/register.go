package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/v3/response"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"github.com/pufferpanel/pufferpanel/v2/web/handlers"
	"gopkg.in/go-playground/validator.v9"
)

func RegisterPost(c *gin.Context) {
	res := response.From(c)
	db := handlers.GetDatabase(c)
	us := &services.User{DB: db}

	res.Fail()
	res.Message("unknown error occurred")

	request := &registerRequest{}
	err := c.BindJSON(request)

	if err != nil {
		res.Fail().Status(400).Data(err)
		return
	}

	validate := validator.New()
	err = validate.Struct(request.Data)
	if err != nil {
		res.Fail().Status(400).Data(err)
		return
	}

	user := &models.User{Username: request.Data.Username, Email: request.Data.Email}
	err = user.SetPassword(request.Data.Password)
	if err != nil {
		res.Fail().Status(400).Error(err)
		return
	}

	err = us.Create(user)
	if err != nil {
		res.Fail().Status(400).Error(err)
		return
	}

	res.Success().Message("")
}

type registerRequest struct {
	Data registerRequestData `json:"data"`
}

type registerRequestData struct {
	Username string `json:"username" validate:"min=3,printascii,required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
