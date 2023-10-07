package tests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestServers(t *testing.T) {
	t.Run("CreateServer", func(t *testing.T) {
		session, err := createSessionAdmin()
		if !assert.NoError(t, err) {
			return
		}

		response := CallAPIRaw("PUT", "/api/servers/testserver", CreateServerData, session)
		assert.Equal(t, http.StatusOK, response.Code)
	})
}
