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

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/middleware/handlers"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func registerSelf(g *gin.RouterGroup) {
	g.Handle("GET", "", handlers.OAuth2Handler(pufferpanel.ScopeNone, false), getSelf)
	g.Handle("PUT", "", handlers.OAuth2Handler(pufferpanel.ScopeNone, false), updateSelf)
	g.Handle("OPTIONS", "", response.CreateOptions("GET", "PUT"))

	g.Handle("GET", "/otp", handlers.OAuth2Handler(pufferpanel.ScopeNone, false), getOtpStatus)
	g.Handle("POST", "/otp", handlers.OAuth2Handler(pufferpanel.ScopeNone, false), startOtpEnroll)
	g.Handle("PUT", "/otp", handlers.OAuth2Handler(pufferpanel.ScopeNone, false), validateOtpEnroll)
	g.Handle("OPTIONS", "/otp", response.CreateOptions("GET", "POST", "PUT"))

	g.Handle("DELETE", "/otp/:token", handlers.OAuth2Handler(pufferpanel.ScopeNone, false), disableOtp)
	g.Handle("OPTIONS", "/otp/:token", response.CreateOptions("DELETE"))

	g.Handle("GET", "/oauth2", handlers.OAuth2Handler(pufferpanel.ScopeNone, false), getPersonalOAuth2Clients)
	g.Handle("POST", "/oauth2", handlers.OAuth2Handler(pufferpanel.ScopeNone, false), createPersonalOAuth2Client)
	g.Handle("OPTIONS", "/oauth2", response.CreateOptions("GET", "POST"))

	g.Handle("DELETE", "/oauth2/:clientId", handlers.OAuth2Handler(pufferpanel.ScopeNone, false), deletePersonalOAuth2Client)
	g.Handle("OPTIONS", "/oauth2/:clientId", response.CreateOptions("DELETE"))
}

// @Summary Get your user info
// @Description Gets the user information of the current user
// @Accept json
// @Produce json
// @Success 200 {object} models.UserView
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /api/self [get]
func getSelf(c *gin.Context) {
	t, exist := c.Get("user")
	user, ok := t.(*models.User)

	if !exist || !ok {
		response.HandleError(c, pufferpanel.ErrUnknownError, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, models.FromUser(user))
}

// @Summary Update your user
// @Description Update user information for your current user
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param user body models.User true "User information"
// @Router /api/self [PUT]
func updateSelf(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	t, exist := c.Get("user")
	user, ok := t.(*models.User)

	if !exist || !ok {
		response.HandleError(c, pufferpanel.ErrUnknownError, http.StatusInternalServerError)
		return
	}

	var viewModel models.UserView
	if err := c.BindJSON(&viewModel); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	if err := viewModel.Valid(true); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	if viewModel.Password == "" {
		response.HandleError(c, pufferpanel.ErrFieldRequired("password"), http.StatusBadRequest)
		return
	}

	if !us.IsValidCredentials(user, viewModel.Password) {
		response.HandleError(c, pufferpanel.ErrInvalidCredentials, http.StatusInternalServerError)
		return
	}

	var oldEmail string
	if user.Email != viewModel.Email {
		oldEmail = user.Email
	}

	viewModel.CopyToModel(user)

	passwordChanged := false
	if viewModel.NewPassword != "" {
		passwordChanged = true
		err := user.SetPassword(viewModel.NewPassword)
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}
	}

	if err := us.Update(user); response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	if oldEmail != "" {
		err := services.GetEmailService().SendEmail(oldEmail, "emailChanged", map[string]interface{}{
			"NEW_EMAIL": user.Email,
		}, true)
		if err != nil {
			logging.Error.Printf("Error sending email: %s\n", err)
		}
	}

	if passwordChanged {
		err := services.GetEmailService().SendEmail(user.Email, "passwordChanged", nil, true)
		if err != nil {
			logging.Error.Printf("Error sending email: %s\n", err)
		}
	}

	c.Status(http.StatusNoContent)
}

func getOtpStatus(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	t, exist := c.Get("user")
	user, ok := t.(*models.User)

	if !exist || !ok {
		response.HandleError(c, pufferpanel.ErrUnknownError, http.StatusInternalServerError)
		return
	}

	otpEnabled, err := us.GetOtpStatus(user.ID)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"otpEnabled": otpEnabled,
	})
}

func startOtpEnroll(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	t, exist := c.Get("user")
	user, ok := t.(*models.User)

	if !exist || !ok {
		response.HandleError(c, pufferpanel.ErrUnknownError, http.StatusInternalServerError)
		return
	}

	secret, img, err := us.StartOtpEnroll(user.ID)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"secret": secret,
		"img":    img,
	})
}

func validateOtpEnroll(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	t, exist := c.Get("user")
	user, ok := t.(*models.User)

	if !exist || !ok {
		response.HandleError(c, pufferpanel.ErrUnknownError, http.StatusInternalServerError)
		return
	}

	request := &ValidateOtpRequest{}

	err := c.BindJSON(request)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	err = us.ValidateOtpEnroll(user.ID, request.Token)
	if err == pufferpanel.ErrInvalidCredentials {
		response.HandleError(c, err, http.StatusBadRequest)
		return
	}
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = services.GetEmailService().SendEmail(user.Email, "otpEnabled", nil, true)
	if err != nil {
		logging.Error.Printf("Error sending email: %s\n", err)
	}
	c.Status(http.StatusNoContent)
}

func disableOtp(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	t, exist := c.Get("user")
	user, ok := t.(*models.User)

	if !exist || !ok {
		response.HandleError(c, pufferpanel.ErrUnknownError, http.StatusInternalServerError)
		return
	}

	err := us.DisableOtp(user.ID, c.Param("token"))
	if err == pufferpanel.ErrInvalidCredentials {
		response.HandleError(c, err, http.StatusBadRequest)
		return
	}
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = services.GetEmailService().SendEmail(user.Email, "otpDisabled", nil, true)
	if err != nil {
		logging.Error.Printf("Error sending email: %s\n", err)
	}
	c.Status(http.StatusNoContent)
}

// @Summary Gets registered oauth2 clients under this user
// @Accept json
// @Produce json
// @Success 200 {object} []models.Client
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /api/self/oauth2 [GET]
func getPersonalOAuth2Clients(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	db := middleware.GetDatabase(c)
	os := &services.OAuth2{DB: db}

	clients, err := os.GetForUserAndServer(user.ID, "")
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	data := make([]*models.Client, 0)
	for _, v := range clients {
		if v.ServerId == "" {
			data = append(data, v)
		}
	}

	c.JSON(http.StatusOK, data)
}

// @Summary Create an account-level OAuth2 client
// @Accept json
// @Produce json
// @Success 200 {object} models.CreatedClient
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param client body models.Client false "Information for the client to create"
// @Router /api/self/oauth2 [POST]
func createPersonalOAuth2Client(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	db := middleware.GetDatabase(c)
	os := &services.OAuth2{DB: db}

	var request models.Client
	err := c.BindJSON(&request)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	client := &models.Client{
		ClientId:    uuid.NewV4().String(),
		UserId:      user.ID,
		Name:        request.Name,
		Description: request.Description,
	}

	secret, err := pufferpanel.GenerateRandomString(36)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = client.SetClientSecret(secret)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = os.Update(client)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = services.GetEmailService().SendEmail(user.Email, "oauthCreated", nil, true)
	if err != nil {
		logging.Error.Printf("Error sending email: %s\n", err)
	}

	c.JSON(http.StatusOK, models.CreatedClient{
		ClientId:     client.ClientId,
		ClientSecret: secret,
	})
}

// @Summary Deletes an account-level OAuth2 client
// @Accept json
// @Produce json
// @Success 204 {object} response.Empty
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Param id path string true "Information for the client to create"
// @Router /api/self/oauth2/{id} [DELETE]
func deletePersonalOAuth2Client(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	clientId := c.Param("clientId")

	db := middleware.GetDatabase(c)
	os := &services.OAuth2{DB: db}

	client, err := os.Get(clientId)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	//ensure the client id is specific for this server, and this user
	if client.UserId != user.ID || client.ServerId != "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	err = os.Delete(client)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = services.GetEmailService().SendEmail(user.Email, "oauthDeleted", nil, true)
	if err != nil {
		logging.Error.Printf("Error sending email: %s\n", err)
	}
	c.Status(http.StatusNoContent)
}

type ValidateOtpRequest struct {
	Token string `json:"token"`
}
