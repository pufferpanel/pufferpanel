package tests

import (
	"encoding/json"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestOauth(t *testing.T) {
	var clientId string
	name := "Test Client"
	description := "this is a test to make sure things even work"

	session, err := createSessionAdmin()
	if !assert.NoError(t, err) {
		return
	}

	t.Run("CreateClient", func(t *testing.T) {
		response := CallAPI("POST", "/api/self/oauth2", map[string]string{
			"name":        name,
			"description": description,
		}, session)
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}
		var client models.Client
		err = json.NewDecoder(response.Body).Decode(&client)
		if !assert.NoError(t, err) {
			return
		}
		clientId = client.ClientId
	})

	t.Run("GetClient", func(t *testing.T) {
		response := CallAPI("GET", "/api/self/oauth2", nil, session)
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}

		var clients []models.Client
		err = json.NewDecoder(response.Body).Decode(&clients)
		if !assert.NoError(t, err) {
			return
		}

		if !assert.NotEmpty(t, clients) {
			return
		}

		found := false
		for _, v := range clients {
			if v.ClientId == clientId {
				found = true
				if !assert.Equal(t, name, v.Name) {
					return
				}
				if !assert.Equal(t, description, v.Description) {
					return
				}
			}
		}
		assert.True(t, found)
	})
}
