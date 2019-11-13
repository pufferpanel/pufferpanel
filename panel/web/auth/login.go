package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v2/panel/services"
	"github.com/pufferpanel/pufferpanel/v2/panel/web/handlers"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"github.com/pufferpanel/pufferpanel/v2/scope"
	"net/http"
)

func LoginPost(c *gin.Context) {
	db := handlers.GetDatabase(c)
	us := &services.User{DB: db}
	ps := &services.Permission{DB: db}

	request := &LoginRequestData{}

	err := c.BindJSON(request)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	user, session, err := us.Login(request.Email, request.Password)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	perms, err := ps.GetForUserAndServer(user.ID, nil)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	data := &LoginResponse{}
	data.Session = session
	data.Scopes = perms.ToScopes()

	c.JSON(http.StatusOK, data)
}

type LoginRequestData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Session string        `json:"session"`
	Scopes  []scope.Scope `json:"scopes,omitempty"`
}
