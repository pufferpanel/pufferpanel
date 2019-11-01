package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/v4/logging"
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

	request := &registerRequestData{}
	err := c.BindJSON(request)

	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	user := &models.User{Username: request.Username, Email: request.Email}
	err = user.SetPassword(request.Password)
	if response.HandleError(c, err, http.StatusInternalServerError) {
	}

	err = us.Create(user)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	//TODO: Have this be an optional flag
	token := ""
	if true {
		_, token, err = us.Login(user.Email, request.Password)
		if err != nil {
			logging.Exception("Error trying to auto-login after register", err)
			c.JSON(200, &registerResponse{Success: true})
			return
		}
	}

	c.JSON(200, &registerResponse{Success: true, Token: token})
}

type registerResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
}

type registerRequestData struct {
	Username string `json:"username" validate:"min=3,printascii,required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
