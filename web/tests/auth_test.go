package tests

import (
	"encoding/json"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/web/auth"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	t.Run("GoodLoginButNoScope", func(t *testing.T) {
		response := CallAPI("POST", "/auth/login", auth.LoginRequestData{
			Email:    loginNoLoginUser.Email,
			Password: loginNoLoginUserPassword,
		}, false)
		assert.Equal(t, http.StatusForbidden, response.Code)
	})
	t.Run("GoodLoginWithLoginScope", func(t *testing.T) {
		response := CallAPI("POST", "/auth/login", auth.LoginRequestData{
			Email:    loginNoServerViewUser.Email,
			Password: loginNoServerViewUserPassword,
		}, false)
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}
		res := &auth.LoginResponse{}
		err := json.NewDecoder(response.Body).Decode(res)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, []*pufferpanel.Scope{pufferpanel.ScopeLogin}, res.Scopes)
	})
	t.Run("GoodLoginWithAdminScope", func(t *testing.T) {
		response := CallAPI("POST", "/auth/login", auth.LoginRequestData{
			Email:    loginAdminUser.Email,
			Password: loginAdminUserPassword,
		}, false)
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}
		res := &auth.LoginResponse{}
		err := json.NewDecoder(response.Body).Decode(res)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, []*pufferpanel.Scope{pufferpanel.ScopeAdmin}, res.Scopes)
	})
	t.Run("NoDataLogin", func(t *testing.T) {
		response := CallAPI("POST", "/auth/login", auth.LoginRequestData{
			Email:    "",
			Password: "",
		}, false)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
	t.Run("InvalidEmail", func(t *testing.T) {
		response := CallAPI("POST", "/auth/login", auth.LoginRequestData{
			Email:    "test@notreal.com",
			Password: "testing123",
		}, false)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
	t.Run("InvalidPassword", func(t *testing.T) {
		response := CallAPI("POST", "/auth/login", auth.LoginRequestData{
			Email:    "test@example.com",
			Password: "testing",
		}, false)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}
