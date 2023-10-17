package tests

import (
	"encoding/json"
	"github.com/pufferpanel/pufferpanel/v3"
	models "github.com/pufferpanel/pufferpanel/v3/oauth2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestOAuth(t *testing.T) {
	var accessToken string

	t.Run("testAuth", func(t *testing.T) {
		form := url.Values{}
		form.Set("grant_type", "client_credentials")
		form.Set("client_id", loginOAuth2Admin.ClientId)
		form.Set("client_secret", loginOAuth2AdminSecret)

		request, _ := http.NewRequest("POST", "/oauth2/token", strings.NewReader(form.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		pufferpanel.Engine.ServeHTTP(response, request)

		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}
		var res models.TokenResponse
		err := json.NewDecoder(response.Body).Decode(&res)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NotEmpty(t, res.AccessToken) {
			return
		}
		accessToken = res.AccessToken
	})

	t.Run("testSelf", func(t *testing.T) {
		response := CallAPI("GET", "/api/self", nil, accessToken)
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}
	})
}
