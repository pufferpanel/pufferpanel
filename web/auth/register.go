/*
 Copyright 2020 Padduck, LLC
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at
  	http://www.apache.org/licenses/LICENSE-2.0
  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/middleware"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

func RegisterPost(c *gin.Context) {
	db := middleware.GetDatabase(c)
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

	ps := &services.Permission{DB: db}
	perms, err := ps.GetForUserAndServer(user.ID, nil)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	perms.ViewServer = true
	perms.Login = true

	err = ps.UpdatePermissions(perms)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	//TODO: Have this be an optional flag
	token := ""
	if true {
		err = services.GetEmailService().SendEmail(user.Email, "accountCreation", nil, true)
		if err != nil {
			logging.Error().Printf("Error sending email: %s", err.Error())
		}

		_, token, err = us.Login(user.Email, request.Password)
		if err != nil {
			logging.Error().Printf("Error trying to auto-login after register: %s", err.Error())
		}
	} else {
		//TODO: Send an email to tell them to validate email
		_ = services.GetEmailService().SendEmail(user.Email, "accountCreation", nil, true)
		if err != nil {
			logging.Error().Printf("Error sending email: %s", err.Error())
		}
	}

	c.JSON(200, &registerResponse{Success: true, Token: token})
}

type registerResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
}

type registerRequestData struct {
	Username string `json:"username" validate:"required,printascii,min=5,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
