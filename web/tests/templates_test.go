package tests

import (
	"encoding/json"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestTemplateAPI(t *testing.T) {
	t.Run("getRepos", func(t *testing.T) {
		session, err := createSessionAdmin()
		if !assert.NoError(t, err) {
			return
		}

		response := CallAPI("GET", "/api/templates", nil, session)
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}
		var templates []*models.TemplateRepo
		err = json.NewDecoder(response.Body).Decode(&templates)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NotEmpty(t, templates) {
			return
		}
		hasLocal := false
		hasCommunity := false
		for _, v := range templates {
			if v.IsLocal {
				hasLocal = true
			} else if v.Name == "community" {
				hasCommunity = true
			}
		}

		assert.True(t, hasLocal, "No local repo")
		assert.True(t, hasCommunity, "No community template repo")
	})

	t.Run("getCommunityRepo", func(t *testing.T) {
		session, err := createSessionAdmin()
		if !assert.NoError(t, err) {
			return
		}

		response := CallAPI("GET", "/api/templates/1", nil, session)
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}
		var templates []*models.Template
		err = json.NewDecoder(response.Body).Decode(&templates)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NotEmpty(t, templates) {
			return
		}
	})
}
