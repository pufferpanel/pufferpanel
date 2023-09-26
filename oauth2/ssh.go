package oauth2

import (
	"encoding/json"
	"errors"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"golang.org/x/crypto/ssh"
	"io"
	"net/http"
	"net/url"
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

	request := createRequest(data)

	response, err := pufferpanel.Http().Do(request)
	defer pufferpanel.CloseResponse(response)
	if err != nil {
		logging.Error.Printf("error talking to auth server: %s", err)
		return nil, errors.New("invalid response from authorization server")
	}

	//we should only get a 200, if we get any others, we have a problem
	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			if recurse && RefreshToken() {
				pufferpanel.CloseResponse(response)
				return validateSSH(username, password, false)
			}
		}

		msg, _ := io.ReadAll(response.Body)

		logging.Error.Printf("Error talking to auth server: [%d] [%s]", response.StatusCode, msg)
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
	for _, v := range scopes {

		t := strings.Split(v, ":")
		if len(t) != 2 {
			continue
		}
		serverId := t[0]
		scope := t[1]

		if pufferpanel.ScopeServerSftp.Is(scope) {
			sshPerms.Extensions = make(map[string]string)
			sshPerms.Extensions["server_id"] = serverId
			return sshPerms, nil
		}
	}
	return nil, errors.New("incorrect username or password")
}
