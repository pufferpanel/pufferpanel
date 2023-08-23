package oauth2

import (
	"encoding/json"
	"github.com/pufferpanel/pufferpanel/v3"
	"net/url"
)

func GetInfo(token string, hint string) (TokenInfoResponse, error) {
	data := url.Values{}
	data.Set("token", token)

	if hint != "" {
		data.Set("token_type_hint", hint)
	}

	var info TokenInfoResponse

	request := createRequest(data)
	response, err := pufferpanel.Http().Do(request)
	defer pufferpanel.CloseResponse(response)
	if err != nil {
		return info, err
	}

	err = json.NewDecoder(response.Body).Decode(&info)
	return info, err
}
