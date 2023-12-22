package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pufferpanel/pufferpanel/v3"
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
	db, err := database.GetConnection()
	if !assert.NoError(t, err) {
		return
	}

	session, err := createSessionAdmin()
	if !assert.NoError(t, err) {
		return
	}

	t.Run("CreateServer", func(t *testing.T) {
		response := CallAPIRaw("PUT", "/api/servers/"+serverId, CreateServerData, session)
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}

		var count int64
		err := db.Model(&models.Server{}).Where(&models.Server{Identifier: serverId}).Count(&count).Error
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, int64(1), count)

		if !assert.DirExists(t, filepath.Join(config.ServersFolder.Value(), serverId)) {
			return
		}
	})

	//previous test is a block,so we can now open up a websocket connection and start playing with it
	//the test here is... do we get all 3 types of messages
	statsReceived := false
	messageReceived := false
	statusReceived := false

	addr := fmt.Sprintf("%s:%d", models.LocalNode.PrivateHost, models.LocalNode.PrivatePort)

	u := fmt.Sprintf("ws://%s/api/servers/%s/socket", addr, serverId)

	header := http.Header{}
	header.Set("Authorization", "Bearer "+session)

	c, _, err := websocket.DefaultDialer.Dial(u, header)
	if !assert.NoError(t, err) {
		return
	}
	listening := true
	defer c.Close()

	go func(conn *websocket.Conn) {
		for listening {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("Error on websocket: %s\n", err.Error())
				continue
			}
			if messageType != websocket.TextMessage {
				fmt.Printf("Unexpected message type [%d]: %s\n", messageType, data)
				continue
			}
			var msg map[string]interface{}
			err = json.NewDecoder(bytes.NewReader(data)).Decode(&msg)
			if err != nil {
				fmt.Printf("Failed to decode message: %s\n", err.Error())
				continue
			}

			msgData := msg["data"]

			switch msg["type"].(string) {
			case messages.Console{}.Key():
				var ms messages.Console
				err = pufferpanel.UnmarshalTo(msgData, &ms)
				if err != nil {
					fmt.Printf("Failed to decode message: %s\n", err.Error())
					continue
				}

				if config.ConsoleForward.Value() {
					fmt.Printf("[WEBSOCKET] %s\n", ms.Logs)
				}

				messageReceived = true
			case messages.Status{}.Key():
				statusReceived = true
			case messages.Stat{}.Key():
				statsReceived = true
			default:
				fmt.Printf("unknown message type: %s\n", msg["type"])
				continue
			}
		}
	}(c)

	t.Run("AddSubUser", func(t *testing.T) {
		var data = []byte(`{"scopes": ["server.view", "server.data.view"]}`)
		response := CallAPIRaw("PUT", "/api/servers/"+serverId+"/user/"+loginNoLoginUser.Email, data, session)
		if !assert.Equal(t, http.StatusNoContent, response.Code) {
			return
		}
	})

	t.Run("GetSubUsers", func(t *testing.T) {
		response := CallAPIRaw("GET", "/api/servers/"+serverId+"/user", nil, session)
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}
		//TODO: Check to make sure our user above was added
		var data []*models.UserPermissionsView
		err = json.NewDecoder(response.Body).Decode(&data)
		if !assert.NoError(t, err) {
			return
		}

		if assert.NotEmpty(t, data) {
			return
		}
		found := false
		for _, v := range data {
			if v.Email == loginNoLoginUser.Email {
				var expectedScopes = []*pufferpanel.Scope{
					pufferpanel.ScopeServerView, pufferpanel.ScopeServerViewData,
				}
				if !assert.Equal(t, expectedScopes, v.Scopes) {
					return
				}
				found = true
			}
		}

		if !found {
			assert.Fail(t, "Failed to locate user")
		}
	})

	t.Run("UpdateVariable", func(t *testing.T) {
		motd := "This is a changed MOTD"
		var variables = []byte(`{"motd": "` + motd + `" }`)
		response := CallAPIRaw("POST", "/api/servers/"+serverId+"/data", variables, session)
		if !assert.Equal(t, http.StatusNoContent, response.Code) {
			return
		}

		response = CallAPI("GET", "/api/servers/"+serverId+"/data", variables, session)
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}

		var res map[string]map[string]pufferpanel.Variable
		err := json.NewDecoder(response.Body).Decode(&res)
		if !assert.NoError(t, err) {
			return
		}
		data := res["data"]
		if !assert.Len(t, data, 1) {
			return
		}

		memVar := data["motd"]
		assert.Equal(t, motd, memVar.Value)
	})

	t.Run("GetStats", func(t *testing.T) {
		response := CallAPI("GET", "/api/servers/"+serverId+"/stats", nil, session)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("SendStatsForServers", func(t *testing.T) {
		servers.SendStatsForServers()
	})

	t.Run("GetEmptyFiles", func(t *testing.T) {
		response := CallAPI("GET", "/api/servers/"+serverId+"/file/", nil, session)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("InstallServer", func(t *testing.T) {
		response := CallAPI("POST", "/api/servers/"+serverId+"/install", nil, session)
		if !assert.Equal(t, http.StatusAccepted, response.Code) {
			return
		}

		time.Sleep(100 * time.Millisecond)

		//we expect it to take more than 100ms, so ensure there is an install occurring
		response = CallAPI("GET", "/api/servers/"+serverId+"/status", nil, session)
		assert.Equal(t, http.StatusOK, response.Code)
		var status messages.Status
		err := json.NewDecoder(response.Body).Decode(&status)
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
		response := CallAPI("POST", "/api/servers/"+serverId+"/start", nil, session)
		assert.Equal(t, http.StatusAccepted, response.Code)

		time.Sleep(1000 * time.Millisecond)
	})

	t.Run("StopServer", func(t *testing.T) {
		response := CallAPI("POST", "/api/servers/"+serverId+"/stop", nil, session)
		if !assert.Equal(t, http.StatusAccepted, response.Code) {
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
			if status.Running {
				counter++
			} else {
				break
			}
		}
		if counter >= timeout {
			assert.Fail(t, "Server took too long to stop, assuming test failed")
		}
	})

	listening = false
	_ = c.Close()

	t.Run("Delete", func(t *testing.T) {
		response := CallAPIRaw("DELETE", "/api/servers/"+serverId, nil, session)
		if !assert.Equal(t, http.StatusNoContent, response.Code) {
			return
		}

		//ensure was actually removed
		if !assert.NoDirExists(t, filepath.Join(config.ServersFolder.Value(), serverId)) {
			return
		}

		var count int64
		err := db.Model(&models.Server{}).Where(&models.Server{Identifier: serverId}).Count(&count).Error
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, int64(0), count)
	})

	t.Run("WebSocketReceivedAll", func(t *testing.T) {
		assert.True(t, statsReceived, "Stats were not received")
		assert.True(t, statusReceived, "Status was not received")
		assert.True(t, messageReceived, "Console messages were not received")
	})
}
