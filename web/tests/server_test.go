package tests

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/sftp"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/database"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/oauth2"
	"github.com/pufferpanel/pufferpanel/v3/scopes"
	"github.com/pufferpanel/pufferpanel/v3/servers"
	"github.com/pufferpanel/pufferpanel/v3/services"
	pufferSftp "github.com/pufferpanel/pufferpanel/v3/sftp"
	"github.com/pufferpanel/pufferpanel/v3/utils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestServers(t *testing.T) {
	db, err := database.GetConnection()
	if !assert.NoError(t, err) {
		return
	}

	session, err := createSessionAdmin()
	if !assert.NoError(t, err) {
		return
	}

	type testLocation struct {
		SFTPAuth pufferpanel.SFTPAuthorization
		Node     *models.Node
	}

	RemoteNode.PublicPort = models.LocalNode.PublicPort
	RemoteNode.PrivatePort = models.LocalNode.PrivatePort
	RemoteNode.SFTPPort = models.LocalNode.SFTPPort

	err = db.Create(RemoteNode).Error
	if !assert.NoError(t, err) {
		return
	}

	tests := map[string]testLocation{
		"Local": {
			SFTPAuth: &services.DatabaseSFTPAuthorization{},
			Node:     models.LocalNode,
		},
		"Remote": {
			SFTPAuth: &oauth2.WebSSHAuthorization{},
			Node:     RemoteNode,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			pufferSftp.SetAuthorization(test.SFTPAuth)
			_ = config.ClientSecret.Set(test.Node.Secret, false)
			var ServerId = "testserver-" + strings.ToLower(name)
			t.Run("CreateServer", func(t *testing.T) {
				data := CreateServerData
				data = strings.Replace(data, "{{{INSERTNODEID}}}", fmt.Sprintf("%d", test.Node.ID), 1)
				response := CallAPIRaw("PUT", "/api/servers/"+ServerId, []byte(data), session)
				if !assert.Equal(t, http.StatusOK, response.Code) {
					return
				}

				var count int64
				err := db.Model(&models.Server{}).Where(&models.Server{Identifier: ServerId}).Count(&count).Error
				if !assert.NoError(t, err) {
					return
				}

				if !assert.Equal(t, int64(1), count) {
					return
				}

				if !assert.DirExists(t, filepath.Join(config.ServersFolder.Value(), ServerId)) {
					return
				}

				err = db.Model(&models.Node{}).Count(&count).Error
				if !assert.NoError(t, err) {
					return
				}
				if !assert.Equal(t, int64(1), count) {
					return
				}
			})

			t.Run("EnsureServerListContains1", func(t *testing.T) {
				response := CallAPI("GET", "/api/servers", nil, session)
				if !assert.Equal(t, http.StatusOK, response.Code) {
					return
				}

				var s *models.ServerSearchResponse
				err := json.NewDecoder(response.Body).Decode(&s)
				if !assert.NoError(t, err) {
					return
				}

				if !assert.NotEmpty(t, s) {
					return
				}

				serverExists := false
				for _, v := range s.Servers {
					if v.Identifier == ServerId {
						serverExists = true
					}
				}
				if !assert.True(t, serverExists, "server does not exist in API") {
					return
				}
			})

			t.Run("EnsureSecondUserCannotSee", func(t *testing.T) {
				sess, err := createSession(db, loginDifferentServerUser)
				if !assert.NoError(t, err) {
					return
				}
				response := CallAPI("GET", "/api/servers", nil, sess)
				if !assert.Equal(t, http.StatusOK, response.Code) {
					return
				}

				var s *models.ServerSearchResponse
				err = json.NewDecoder(response.Body).Decode(&s)
				if !assert.NoError(t, err) {
					return
				}

				if !assert.Empty(t, s.Servers) {
					return
				}
			})

			t.Run("AdminUpdate", func(t *testing.T) {
				response := CallAPIRaw("PUT", "/api/servers/"+ServerId+"/definition", EditServerData, session)
				if !assert.Equal(t, http.StatusNoContent, response.Code) {
					return
				}

				var server *models.Server
				err := db.Model(&server).Where(&models.Server{Identifier: ServerId}).Find(&server).Error
				if !assert.NoError(t, err) {
					return
				}

				if !assert.Equal(t, EditServerNewName, server.Name) {
					return
				}
				if !assert.Equal(t, EditServerNewIP, server.IP) {
					return
				}
				if !assert.Equal(t, EditServerNewPort, server.Port) {
					return
				}

				var count int64
				err = db.Model(&models.Node{}).Count(&count).Error
				if !assert.NoError(t, err) {
					return
				}
				if !assert.Equal(t, int64(1), count) {
					return
				}
			})

			t.Run("AdminDataUpdate", func(t *testing.T) {
				response := CallAPIRaw("PUT", "/api/servers/"+ServerId+"/data", NewVariableChanges, session)
				if !assert.Equal(t, http.StatusNoContent, response.Code) {
					return
				}

				var server *models.Server
				err := db.Model(&server).Where(&models.Server{Identifier: ServerId}).Find(&server).Error
				if !assert.NoError(t, err) {
					return
				}

				if !assert.Equal(t, NewVariableChangeIP, server.IP) {
					return
				}
				if !assert.Equal(t, NewVariableChangePort, server.Port) {
					return
				}

				var count int64
				err = db.Model(&models.Node{}).Count(&count).Error
				if !assert.NoError(t, err) {
					return
				}
				if !assert.Equal(t, int64(1), count) {
					return
				}
			})

			if t.Failed() {
				return
			}

			//previous test is a block,so we can now open up a websocket connection and start playing with it
			//the test here is... do we get all 3 types of messages
			statsReceived := false
			messageReceived := false
			statusReceived := false

			addr := fmt.Sprintf("%s:%d", models.LocalNode.PrivateHost, models.LocalNode.PrivatePort)

			u := fmt.Sprintf("ws://%s/api/servers/%s/socket", addr, ServerId)

			header := http.Header{}
			header.Set("Authorization", "Bearer "+session)

			c, _, err := websocket.DefaultDialer.Dial(u, header)
			if !assert.NoError(t, err) {
				return
			}
			listening := true
			defer utils.Close(c)

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
					case pufferpanel.MessageTypeLog:
						var ms pufferpanel.ServerLogs
						err = utils.UnmarshalTo(msgData, &ms)
						if err != nil {
							fmt.Printf("Failed to decode message: %s\n", err.Error())
							continue
						}

						if config.ConsoleForward.Value() {
							fmt.Printf("[WEBSOCKET] %s\n", ms.Logs)
						}

						messageReceived = true
					case pufferpanel.MessageTypeStatus:
						statusReceived = true
					case pufferpanel.MessageTypeStats:
						statsReceived = true
					default:
						fmt.Printf("unknown message type: %s\n", msg["type"])
						continue
					}
				}
			}(c)

			t.Run("AddSubUser", func(t *testing.T) {
				var data = []byte(`{"scopes": ["server.view", "server.data.view"]}`)
				response := CallAPIRaw("PUT", "/api/servers/"+ServerId+"/user/"+loginNoLoginUser.Email, data, session)
				if !assert.Equal(t, http.StatusNoContent, response.Code) {
					return
				}
			})

			t.Run("GetSubUsers", func(t *testing.T) {
				response := CallAPIRaw("GET", "/api/servers/"+ServerId+"/user", nil, session)
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
						var expectedScopes = []*scopes.Scope{
							scopes.ScopeServerView, scopes.ScopeServerViewData,
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
				response := CallAPIRaw("POST", "/api/servers/"+ServerId+"/data", variables, session)
				if !assert.Equal(t, http.StatusNoContent, response.Code) {
					return
				}

				response = CallAPI("GET", "/api/servers/"+ServerId+"/data", variables, session)
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
				response := CallAPI("GET", "/api/servers/"+ServerId+"/stats", nil, session)
				assert.Equal(t, http.StatusOK, response.Code)
			})

			t.Run("GetStatsSecondUserFail", func(t *testing.T) {
				sess, err := createSession(db, loginDifferentServerUser)
				if !assert.NoError(t, err) {
					return
				}
				response := CallAPI("GET", "/api/servers/"+ServerId+"/stats", nil, sess)
				assert.Equal(t, http.StatusForbidden, response.Code)
			})

			t.Run("SendStatsForServers", func(t *testing.T) {
				servers.SendStatsForServers()
			})

			t.Run("GetEmptyFiles", func(t *testing.T) {
				response := CallAPI("GET", "/api/servers/"+ServerId+"/file/", nil, session)
				assert.Equal(t, http.StatusOK, response.Code)
			})

			t.Run("InstallServer", func(t *testing.T) {
				response := CallAPI("POST", "/api/servers/"+ServerId+"/install", nil, session)
				if !assert.Equal(t, http.StatusAccepted, response.Code) {
					return
				}

				time.Sleep(100 * time.Millisecond)

				//we expect it to take more than 100ms, so ensure there is an install occurring
				response = CallAPI("GET", "/api/servers/"+ServerId+"/status", nil, session)
				assert.Equal(t, http.StatusOK, response.Code)
				var msg pufferpanel.ServerRunning
				err := json.NewDecoder(response.Body).Decode(&msg)
				if !assert.NoError(t, err) {
					return
				}
				if !assert.True(t, msg.Installing) {
					return
				}

				//now we wait for the install to finish
				timeout := 60
				counter := 0
				for counter < timeout {
					time.Sleep(time.Second)
					response = CallAPI("GET", "/api/servers/"+ServerId+"/status", nil, session)
					assert.Equal(t, http.StatusOK, response.Code)
					err = json.NewDecoder(response.Body).Decode(&msg)
					if !assert.NoError(t, err) {
						return
					}

					if msg.Installing {
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
				response := CallAPI("POST", "/api/servers/"+ServerId+"/start", nil, session)
				assert.Equal(t, http.StatusAccepted, response.Code)

				time.Sleep(1000 * time.Millisecond)
			})

			t.Run("StopServer", func(t *testing.T) {
				response := CallAPI("POST", "/api/servers/"+ServerId+"/stop", nil, session)
				if !assert.Equal(t, http.StatusAccepted, response.Code) {
					return
				}

				//now we wait for the install to finish
				timeout := 60
				counter := 0
				for counter < timeout {
					time.Sleep(time.Second)
					response = CallAPI("GET", "/api/servers/"+ServerId+"/status", nil, session)
					assert.Equal(t, http.StatusOK, response.Code)
					var msg pufferpanel.ServerRunning
					err = json.NewDecoder(response.Body).Decode(&msg)
					if !assert.NoError(t, err) {
						return
					}

					if msg.Running {
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

			//create a fake file that we can use to both

			dir := filepath.Join(config.ServersFolder.Value(), ServerId, "testarchive")
			err = os.Mkdir(dir, 0755)
			if !assert.NoError(t, err) {
				return
			}

			fileLocation := filepath.Join(dir, "file.img")
			tmpFile, err := os.Create(fileLocation)
			if !assert.NoError(t, err) {
				return
			}

			hasher := sha256.New()
			w := io.MultiWriter(tmpFile, hasher)

			_, err = io.CopyN(w, rand.Reader, 1024*1024*1024)
			if !assert.NoError(t, err) {
				return
			}

			_ = tmpFile.Close()
			expectedHash := hasher.Sum(nil)

			//test other functionality
			t.Run("Archive", func(t *testing.T) {
				response := CallAPI("POST", "/api/servers/"+ServerId+"/archive/archive.zip", []string{"testarchive"}, session)
				if !assert.Equal(t, http.StatusNoContent, response.Code) {
					return
				}
				_ = os.RemoveAll(dir)
			})

			t.Run("Extract", func(t *testing.T) {
				response := CallAPI("POST", "/api/servers/"+ServerId+"/extract/archive.zip", nil, session)
				if !assert.Equal(t, http.StatusNoContent, response.Code) {
					return
				}
				if !assert.FileExists(t, fileLocation) {
					return
				}

				f, err := os.Open(fileLocation)
				if !assert.NoError(t, err) {
					return
				}
				defer utils.Close(f)
				h := sha256.New()
				_, err = io.Copy(h, f)
				if !assert.NoError(t, err) {
					return
				}

				if !assert.Equal(t, expectedHash, h.Sum(nil), "File hashes do not match") {
					return
				}
			})

			t.Run("SFTP", func(t *testing.T) {

				t.Run("Admin", func(t *testing.T) {
					sshConfig := &ssh.ClientConfig{
						User: loginAdminUser.Email + "#" + ServerId,
						Auth: []ssh.AuthMethod{
							ssh.Password(loginAdminUserPassword),
						},
						HostKeyCallback: ssh.InsecureIgnoreHostKey(),
					}

					client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", test.Node.PrivateHost, test.Node.SFTPPort), sshConfig)
					defer utils.Close(client)
					if !assert.NoError(t, err) {
						return
					}
					sftpClient, err := sftp.NewClient(client)
					if !assert.NoError(t, err) {
						return
					}
					files, err := sftpClient.ReadDir("/")
					if !assert.NoError(t, err) {
						return
					}
					if !assert.NotEmpty(t, files) {
						return
					}
				})

				t.Run("User", func(t *testing.T) {
					sshConfig := &ssh.ClientConfig{
						User: loginSftpUser.Email + "#" + ServerId,
						Auth: []ssh.AuthMethod{
							ssh.Password(loginSftpUserPassword),
						},
						HostKeyCallback: ssh.InsecureIgnoreHostKey(),
					}
					client, err := ssh.Dial("tcp", "localhost:5657", sshConfig)
					defer utils.Close(client)
					if !assert.NoError(t, err) {
						return
					}
					sftpClient, err := sftp.NewClient(client)
					if !assert.NoError(t, err) {
						return
					}
					files, err := sftpClient.ReadDir("/")
					if !assert.NoError(t, err) {
						return
					}
					if !assert.NotEmpty(t, files) {
						return
					}
				})

				t.Run("NonUser", func(t *testing.T) {
					sshConfig := &ssh.ClientConfig{
						User: loginDifferentServerUser.Email + "#" + ServerId,
						Auth: []ssh.AuthMethod{
							ssh.Password(loginDifferentServerUserPassword),
						},
						HostKeyCallback: ssh.InsecureIgnoreHostKey(),
					}
					client, err := ssh.Dial("tcp", "localhost:5657", sshConfig)
					defer utils.Close(client)
					if !assert.Error(t, err) {
						return
					}
				})
			})

			t.Run("Delete", func(t *testing.T) {
				response := CallAPIRaw("DELETE", "/api/servers/"+ServerId, nil, session)
				if !assert.Equal(t, http.StatusNoContent, response.Code) {
					return
				}

				//ensure was actually removed
				if !assert.NoDirExists(t, filepath.Join(config.ServersFolder.Value(), ServerId)) {
					return
				}

				var count int64
				err := db.Model(&models.Server{}).Where(&models.Server{Identifier: ServerId}).Count(&count).Error
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
		})
	}
}
