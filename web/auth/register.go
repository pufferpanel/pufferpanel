package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/logging"
	builder "github.com/pufferpanel/apufferi/response"
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/pufferpanel/pufferpanel/services"
	"gopkg.in/go-playground/validator.v9"
)

func RegisterPost(c *gin.Context) {
	response := builder.Respond(c)
	response.Fail()
	response.Message("unknown error occurred")
	defer response.Send()

	request := &registerRequest{}
	err := c.BindJSON(request)

	if err != nil {
		response.Fail().Status(400).Data(err)
		return
	}

	validate := validator.New()
	err = validate.Struct(request.Data)
	if err != nil {
		response.Fail().Status(400).Data(err)
		return
	}

	us, err := services.GetUserService()
	if err != nil {
		response.Fail().Status(400).Message("error loading user service").Data(err)
		logging.Build(logging.ERROR).WithMessage("error loading user service").WithError(err).Log()
		return
	}

	user := &models.User{Username: request.Data.Username, Email: request.Data.Email}
	err = user.SetPassword(request.Data.Password)
	if err != nil {
		response.Fail().Status(400).Data(err)
		return
	}

	err = us.Create(user)
	if err != nil {
		response.Fail().Status(400).Data(err)
		return
	}
	response.Success()
}

type registerRequest struct {
	Data registerRequestData `json:"data"`
}

type registerRequestData struct {
	Username string `json:"username" validate:"min=8,printascii,required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
