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
	"encoding/json"
	"errors"
	"github.com/pufferpanel/pufferpanel/v2/daemon/commons"
	"github.com/pufferpanel/pufferpanel/v2/shared/logging"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
)

type WebSSHAuthorization struct {
}

func (ws *WebSSHAuthorization) Validate(username string, password string) (*ssh.Permissions, error) {
	return validateSSH(username, password, true)
}

func validateSSH(username string, password string, recurse bool) (*ssh.Permissions, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", username)
	data.Set("password", password)
	data.Set("scope", "sftp")
	encodedData := data.Encode()

	request := createRequest(encodedData)

	RefreshIfStale()

	atLocker.RLock()
	request.Header.Add("Authorization", "Bearer "+daemonToken)
	atLocker.RUnlock()
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))

	response, err := client.Do(request)
	defer commons.CloseResponse(response)
	if err != nil {
		logging.Exception("error talking to auth server", err)
		return nil, errors.New("invalid response from authorization server")
	}

	//we should only get a 200, if we get any others, we have a problem
	if response.StatusCode != 200 {
		if response.StatusCode == 401 {
			if recurse && RefreshToken() {
				commons.CloseResponse(response)
				return validateSSH(username, password, false)
			}
		}

		msg, _ := ioutil.ReadAll(response.Body)

		logging.Error("Error talking to auth server: [%d] [%s]", response.StatusCode, msg)
		return nil, errors.New("invalid response from authorization server")
	}

	var respArr map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&respArr)
	if err != nil {
		return nil, err
	}
	if respArr["error"] != nil {
		return nil, errors.New("incorrect username or password")
	}
	sshPerms := &ssh.Permissions{}
	scopes := strings.Split(respArr["scope"].(string), " ")
	if len(scopes) != 2 {
		return nil, errors.New("invalid response from authorization server")
	}
	for _, v := range scopes {
		if v != "sftp" {
			sshPerms.Extensions = make(map[string]string)
			sshPerms.Extensions["server_id"] = v
			return sshPerms, nil
		}
	}
	return nil, errors.New("incorrect username or password")
}
