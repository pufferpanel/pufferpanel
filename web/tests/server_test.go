package tests

import (
	"encoding/json"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/database"
	"github.com/pufferpanel/pufferpanel/v3/messages"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/servers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path/filepath"
	"testing"
	"time"
)

func TestServers(t *testing.T) {
	serverId := "testserver"

	t.Run("CreateServer", func(t *testing.T) {
		session, err := createSessionAdmin()
		if !assert.NoError(t, err) {
			return
		}

		response := CallAPIRaw("PUT", "/api/servers/"+serverId, CreateServerData, session)
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}

		db, err := database.GetConnection()
		if !assert.NoError(t, err) {
			return
		}

		var count int64
		err = db.Model(&models.Server{}).Where(&models.Server{Identifier: serverId}).Count(&count).Error
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, int64(1), count)

		if !assert.DirExists(t, filepath.Join(config.ServersFolder.Value(), serverId)) {
			return
		}
	})

	t.Run("GetStats", func(t *testing.T) {
		session, err := createSessionAdmin()
		if !assert.NoError(t, err) {
			return
		}

		response := CallAPI("GET", "/api/servers/"+serverId+"/stats", nil, session)
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

		response := CallAPI("GET", "/api/servers/"+serverId+"/file/", nil, session)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("InstallServer", func(t *testing.T) {
		session, err := createSessionAdmin()
		if !assert.NoError(t, err) {
			return
		}

		response := CallAPI("POST", "/api/servers/"+serverId+"/install", nil, session)
		if !assert.Equal(t, http.StatusAccepted, response.Code) {
			return
		}

		time.Sleep(100 * time.Millisecond)

		//we expect it to take more than 100ms, so ensure there is an install occurring
		response = CallAPI("GET", "/api/servers/"+serverId+"/status", nil, session)
		assert.Equal(t, http.StatusOK, response.Code)
		var status messages.Status
		err = json.NewDecoder(response.Body).Decode(&status)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.True(t, status.Installing) {
			return
		}

		//now we wait for the install to finish
		timeout := 60
		counter := 0
		for counter < timeout {
			time.Sleep(time.Second)
			response = CallAPI("GET", "/api/servers/"+serverId+"/status", nil, session)
			assert.Equal(t, http.StatusOK, response.Code)
			var status messages.Status
			err = json.NewDecoder(response.Body).Decode(&status)
			if !assert.NoError(t, err) {
				return
			}
			if status.Installing {
				counter++
			} else {
				break
			}
		}
		if counter >= timeout {
			assert.Fail(t, "Server took too long to install, assuming test failed")
		}
	})

	t.Run("StartServer", func(t *testing.T) {
		session, err := createSessionAdmin()
		if !assert.NoError(t, err) {
			return
		}

		response := CallAPI("POST", "/api/servers/"+serverId+"/start", nil, session)
		assert.Equal(t, http.StatusAccepted, response.Code)

		time.Sleep(1000 * time.Millisecond)

		//we expect it to take more than 1 second, so ensure there is a started server
		response = CallAPI("GET", "/api/servers/"+serverId+"/status", nil, session)
		assert.Equal(t, http.StatusOK, response.Code)
		var status messages.Status
		err = json.NewDecoder(response.Body).Decode(&status)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.True(t, status.Running) {
			return
		}
	})

	t.Run("Delete", func(t *testing.T) {
		session, err := createSessionAdmin()
		if !assert.NoError(t, err) {
			return
		}

		response := CallAPIRaw("DELETE", "/api/servers/"+serverId, nil, session)
		if !assert.Equal(t, http.StatusNoContent, response.Code) {
			return
		}

		db, err := database.GetConnection()
		if !assert.NoError(t, err) {
			return
		}

		//ensure was actually removed
		if !assert.NoDirExists(t, filepath.Join(config.ServersFolder.Value(), serverId)) {
			return
		}

		var count int64
		err = db.Model(&models.Server{}).Where(&models.Server{Identifier: serverId}).Count(&count).Error
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, int64(0), count)
	})
}
