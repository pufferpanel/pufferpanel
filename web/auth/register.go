package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/v4/response"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"github.com/pufferpanel/pufferpanel/v2/web/handlers"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

func RegisterPost(c *gin.Context) {
	db := handlers.GetDatabase(c)
	us := &services.User{DB: db}

	request := &registerRequest{}
	err := c.BindJSON(request)

	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	validate := validator.New()
	err = validate.Struct(request.Data)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	user := &models.User{Username: request.Data.Username, Email: request.Data.Email}
	err = user.SetPassword(request.Data.Password)
	if response.HandleError(c, err, http.StatusInternalServerError) {
	}

	err = us.Create(user)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}
}

type registerRequest struct {
	Data registerRequestData `json:"data"`
}

type registerRequestData struct {
	Username string `json:"username" validate:"min=3,printascii,required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
