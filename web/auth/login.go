package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"net/http"
	"time"
)

func LoginPost(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	request := &LoginRequestData{}

	err := c.BindJSON(request)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	user, otpNeeded, err := us.ValidateLogin(request.Email, request.Password)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	if otpNeeded {
		userSession := sessions.Default(c)
		userSession.Set("user", user.Email)
		userSession.Set("time", time.Now().Unix())
		_ = userSession.Save()
		c.JSON(http.StatusOK, &LoginResponse{
			OtpNeeded: true,
		})
		return
	}

	createSession(c, user)
}

func OtpPost(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

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

	user, err := us.ValidOtp(email, request.Token)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	createSession(c, user)
}

func createSession(c *gin.Context, user *models.User) {
	db := middleware.GetDatabase(c)
	ps := &services.Permission{DB: db}
	ss := &services.Session{DB: db}

	perms, err := ps.GetForUserAndServer(user.ID, nil)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	if !pufferpanel.ContainsScope(perms.Scopes, pufferpanel.ScopeLogin) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	session, err := ss.CreateForUser(user)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	data := &LoginResponse{}
	data.Scopes = perms.Scopes

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

type LoginResponse struct {
	Scopes    []*pufferpanel.Scope `json:"scopes,omitempty"`
	OtpNeeded bool                 `json:"otpNeeded,omitempty"`
}

type OtpRequestData struct {
	Token string `json:"token"`
}
