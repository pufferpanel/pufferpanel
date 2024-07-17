package oauth2

import (
	"bytes"
	"github.com/gin-gonic/gin/binding"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"net/http"
	"net/url"
)

func createRequest(data url.Values) (request *http.Request) {
	authUrl := config.AuthUrl.Value()
	request, _ = http.NewRequest("POST", authUrl, bytes.NewBufferString(data.Encode()))

	request.Header.Add("Authorization", "Bearer "+config.ClientSecret.Value())
	request.Header.Add("Content-Type", binding.MIMEPOSTForm)
	return
}

type TokenInfoResponse struct {
	Active bool   `json:"active"`
	Scope  string `json:"scope,omitempty"`
	ErrorResponse
} //@name OAuth2TokenInfoResponse

type TokenResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
	Scope       string `json:"scope"`
	ErrorResponse
} //@name OAuth2TokenResponse

type ErrorResponse struct {
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
} //@name OAuthErrorResponse
