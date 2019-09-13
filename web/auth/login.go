package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/v3/response"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/services"
)

func LoginPost(c *gin.Context) {
	res := response.From(c)

	request := &loginRequest{}

	err := c.BindJSON(request)
	if err != nil {
		res.Message("invalid request").Status(400).Error(err).Fail()
		return
	}

	db, err := database.GetConnection()
	if response.HandleError(res, err) {
		return
	}

	us := &services.User{DB: db}
	os := services.GetOAuth(db)

	session, err := us.Login(request.Data.Email, request.Data.Password)
	if response.HandleError(res, err) {
		return
	}

	data := &loginResponse{}
	data.Session = session

	_, client, err := os.GetByToken(data.Session)

	if response.HandleError(res, err) {
		return
	}

	data.ServerScopes = make(map[string][]string)

	for _, v := range client.ServerScopes {
		var serverName string
		if v.ServerId != nil {
			serverName = v.Server.Identifier
		} else {
			serverName = ""
		}
		m := data.ServerScopes[serverName]
		if m == nil {
			m = []string{v.Scope}
		} else {
			m = append(m, v.Scope)
		}
		data.ServerScopes[serverName] = m
	}

	res.Data(data)
}

type loginRequest struct {
	Data loginRequestData `json:"data"`
}

type loginRequestData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Session      string              `json:"session"`
	ServerScopes map[string][]string `json:"scopes"`
}
