package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/logging"
	builder "github.com/pufferpanel/apufferi/response"
	"github.com/pufferpanel/pufferpanel/errors"
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
	if us == nil && err == nil {
		err = errors.ErrServiceNotAvailable
	}
	if err != nil {
		response.Fail().Status(400).Error(err)
		logging.Build(logging.ERROR).WithMessage("error loading user service").WithError(err).Log()
		return
	}

	user := &models.User{Username: request.Data.Username, Email: request.Data.Email}
	err = user.SetPassword(request.Data.Password)
	if err != nil {
		response.Fail().Status(400).Error(err)
		return
	}

	err = us.Create(user)
	if err != nil {
		response.Fail().Status(400).Error(err)
		return
	}

	os, err := services.GetOAuthService()
	if err != nil {
		response.Fail().Status(400).Error(err)
		return
	}

	client, _, err := os.GetByUser(user)
	if err != nil {
		response.Fail().Status(400).Error(err)
		return
	}

	err = os.AddScope(client, nil, "servers.view")
	if err != nil {
		response.Fail().Status(400).Error(err)
		return
	}

	response.Success().Message("")
}

type registerRequest struct {
	Data registerRequestData `json:"data"`
}

type registerRequestData struct {
	Username string `json:"username" validate:"min=3,printascii,required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
