/*
 Copyright 2019 Padduck, LLC

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

package oauth2

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin/binding"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/config"
	"github.com/pufferpanel/pufferpanel/v2/logging"
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
	Active           bool   `json:"active"`
	Scope            string `json:"scope,omitempty"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

type TokenResponse struct {
	AccessToken      string `json:"access_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	ExpiresIn        int64  `json:"expires_in,omitempty"`
	Scope            string `json:"scope"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}
