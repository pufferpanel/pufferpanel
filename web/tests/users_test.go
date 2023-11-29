package tests

import (
	"encoding/json"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestUserApi(t *testing.T) {
	session, err := createSessionAdmin()
	if !assert.NoError(t, err) {
		return
	}

	t.Run("TestUserWithPerms", func(t *testing.T) {
		url := fmt.Sprintf("/api/users/%d/perms", loginNoServerViewUser.ID)
		response := CallAPI("GET", url, nil, session)
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}

		perms := &models.PermissionView{}
		err = json.NewDecoder(response.Body).Decode(perms)
		if !assert.NoError(t, err) {
			return
		}

		if !assert.NotEmpty(t, perms.Scopes) {
			return
		}
	})

	var newUserId uint

	t.Run("CreateUser", func(t *testing.T) {
		response := CallAPI("POST", "/api/users", map[string]string{
			"username": "apicreateduser",
			"email":    "none@example.com",
			"password": "testing123!",
		}, session)
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}

		user := &models.UserView{}
		err = json.NewDecoder(response.Body).Decode(user)
		if !assert.NoError(t, err) {
			return
		}
		newUserId = user.Id
	})

	t.Run("TestEmptyUserPerms", func(t *testing.T) {
		url := fmt.Sprintf("/api/users/%d/perms", newUserId)
		response := CallAPI("GET", url, nil, session)
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}

		perms := &models.PermissionView{}
		err = json.NewDecoder(response.Body).Decode(perms)
		if !assert.NoError(t, err) {
			return
		}

		if !assert.Empty(t, perms.Scopes) {
			return
		}
	})

	t.Run("TestSetEmptyPerms", func(t *testing.T) {
		url := fmt.Sprintf("/api/users/%d/perms", newUserId)
		response := CallAPI("PUT", url, map[string]interface{}{
			"scopes": []string{},
		}, session)
		if !assert.Equal(t, http.StatusNoContent, response.Code) {
			return
		}

		response = CallAPI("GET", url, nil, session)
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}

		perms := &models.PermissionView{}
		err = json.NewDecoder(response.Body).Decode(perms)
		if !assert.NoError(t, err) {
			return
		}

		if !assert.Empty(t, perms.Scopes) {
			return
		}
	})

	t.Run("TestDelete", func(t *testing.T) {
		url := fmt.Sprintf("/api/users/%d", newUserId)
		response := CallAPI("DELETE", url, nil, session)
		if !assert.Equal(t, http.StatusNoContent, response.Code) {
			return
		}
	})
}
