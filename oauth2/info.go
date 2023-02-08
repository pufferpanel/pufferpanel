package oauth2

import (
	"encoding/json"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"net/http"
	"net/url"
)

func GetInfo(token string, hint string) (info TokenInfoResponse, err error) {
	client := pufferpanel.Http()
	request := &http.Request{}

	data := url.Values{}
	data.Set("token", token)

	if hint != "" {
		data.Set("token_type_hint", hint)
	}

	request.URL, err = url.Parse(config.AuthUrl.Value())
	if err != nil {
		return
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", "Bearer ")

	res, err := client.Do(request)

	err = json.NewDecoder(res.Body).Decode(&info)
	return
}
