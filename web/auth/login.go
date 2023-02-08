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
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/middleware/panelmiddleware"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"net/http"
	"time"
)

func LoginPost(c *gin.Context) {
	db := panelmiddleware.GetDatabase(c)
	us := &services.User{DB: db}
	ps := &services.Permission{DB: db}

	request := &LoginRequestData{}

	err := c.BindJSON(request)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	user, session, otpNeeded, err := us.Login(request.Email, request.Password)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	if otpNeeded {
		userSession := sessions.Default(c)
		userSession.Set("user", user.Email)
		userSession.Set("time", time.Now().Unix())
		_ = userSession.Save()
		c.JSON(http.StatusOK, &LoginOtpResponse{
			OtpNeeded: true,
		})
		return
	}

	perms, err := ps.GetForUserAndServer(user.ID, nil)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	data := &LoginResponse{}
	//data.Session = session
	data.Scopes = perms.ToScopes()

	secure := false
	if c.Request.TLS != nil {
		secure = true
	}

	maxAge := int(time.Hour / time.Second)

	c.SetCookie("puffer_auth", session, maxAge, "/", "", secure, true)
	c.SetCookie("puffer_auth_expires", "", maxAge, "/", "", secure, false)

	c.JSON(http.StatusOK, data)
}

func OtpPost(c *gin.Context) {
	db := panelmiddleware.GetDatabase(c)
	us := &services.User{DB: db}
	ps := &services.Permission{DB: db}

	request := &OtpRequestData{}

	err := c.BindJSON(request)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	userSession := sessions.Default(c)
	email := userSession.Get("user").(string)
	timestamp := userSession.Get("time").(int64)

	if email == "" {
		response.HandleError(c, pufferpanel.ErrInvalidSession, http.StatusBadRequest)
		return
	}

	if timestamp < time.Now().Unix()-300 {
		userSession.Clear()
		_ = userSession.Save()
		response.HandleError(c, pufferpanel.ErrSessionExpired, http.StatusBadRequest)
		return
	}

	user, session, err := us.LoginOtp(email, request.Token)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	perms, err := ps.GetForUserAndServer(user.ID, nil)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	data := &LoginResponse{}
	//data.Session = session
	data.Scopes = perms.ToScopes()

	secure := false
	if c.Request.TLS != nil {
		secure = true
	}

	maxAge := int(time.Hour / time.Second)

	c.SetCookie("puffer_auth", session, maxAge, "/", "", secure, true)
	c.SetCookie("puffer_auth_expires", "", maxAge, "/", "", secure, false)

	c.JSON(http.StatusOK, data)
}

type LoginRequestData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOtpResponse struct {
	OtpNeeded bool `json:"otpNeeded"`
}

type LoginResponse struct {
	Scopes []pufferpanel.Scope `json:"scopes,omitempty"`
}

type OtpRequestData struct {
	Token string `json:"token"`
}
