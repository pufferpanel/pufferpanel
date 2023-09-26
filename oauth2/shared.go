package oauth2

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin/binding"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var atLocker = &sync.RWMutex{}
var daemonToken string
var lastRefresh time.Time
var expiresIn int64

func RefreshToken() bool {
	atLocker.Lock()
	defer atLocker.Unlock()

	//if we just refreshed in the last minute, don't refresh the token
	if lastRefresh.Add(1 * time.Minute).After(time.Now()) {
		return false
	}

	clientId := config.ClientId.Value()
	if clientId == "" {
		logging.Error.Printf("error talking to auth server: %s", errors.New("client id not specified"))
		return false
	}

	clientSecret := config.ClientSecret.Value()
	if clientSecret == "" {
		logging.Error.Printf("error talking to auth server: %s", errors.New("client secret not specified"))
		return false
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)

	authUrl := config.AuthUrl.Value()
	request, _ := http.NewRequest("POST", authUrl, bytes.NewBufferString(data.Encode()))
	request.Header.Add("Content-Type", binding.MIMEPOSTForm)

	response, err := pufferpanel.Http().Do(request)
	defer pufferpanel.CloseResponse(response)
	if err != nil {
		logging.Error.Printf("error talking to auth server: %s", err)
		return false
	}

	var responseData TokenResponse
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		logging.Error.Printf("error talking to auth server: %s", err.Error())
		return false
	}

	if responseData.Error != "" {
		logging.Error.Printf("error talking to auth server: %s", responseData.Error)
		return false
	}

	daemonToken = responseData.AccessToken
	lastRefresh = time.Now()
	expiresIn = responseData.ExpiresIn

	return true
}

func RefreshIfStale() {
	//we know the token only lasts about an hour,
	//so we'll check to see if we know the cache is old
	atLocker.RLock()
	oldCache := lastRefresh.Add(time.Second * time.Duration(expiresIn)).Before(time.Now())
	atLocker.RUnlock()
	if oldCache {
		RefreshToken()
	}
}

func createRequest(data url.Values) (request *http.Request) {
	authUrl := config.AuthUrl.Value()
	request, _ = http.NewRequest("POST", authUrl, bytes.NewBufferString(data.Encode()))

	RefreshIfStale()

	atLocker.RLock()
	request.Header.Add("Authorization", "Bearer "+daemonToken)
	atLocker.RUnlock()
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
