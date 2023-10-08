package tests

import (
	"github.com/pufferpanel/pufferpanel/v3/servers"
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

	t.Run("StartServer", func(t *testing.T) {
		session, err := createSessionAdmin()
		if !assert.NoError(t, err) {
			return
		}

		response := CallAPI("POST", "/api/servers/testserver/start", nil, session)
		assert.Equal(t, http.StatusAccepted, response.Code)
	})

	t.Run("GetStats", func(t *testing.T) {
		session, err := createSessionAdmin()
		if !assert.NoError(t, err) {
			return
		}

		response := CallAPI("GET", "/api/servers/testserver/stats", nil, session)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("SendStatsForServers", func(t *testing.T) {
		servers.SendStatsForServers()
	})

	t.Run("GetEmptyFiles", func(t *testing.T) {
		session, err := createSessionAdmin()
		if !assert.NoError(t, err) {
			return
		}

		response := CallAPI("GET", "/api/servers/testserver/file/", nil, session)
		assert.Equal(t, http.StatusOK, response.Code)
	})
}
