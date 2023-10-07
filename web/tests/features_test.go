package tests

import (
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestFeatures(t *testing.T) {
	t.Run("TestFeatures", func(t *testing.T) {
		session, err := createSessionAdmin()
		if !assert.NoError(t, err) {
			return
		}

		url := fmt.Sprintf("/api/nodes/%d/features", models.LocalNode.ID)
		response := CallAPI("GET", url, nil, session)
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}
		if !assert.NotEmpty(t, response.Body.String()) {
			return
		}
	})
}
