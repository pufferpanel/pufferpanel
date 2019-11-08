package server

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pufferpanel/pufferpanel/v2/daemon/messages"
	"github.com/pufferpanel/pufferpanel/v2/daemon/programs"
	"github.com/pufferpanel/pufferpanel/v2/shared"
	"github.com/pufferpanel/pufferpanel/v2/shared/logging"
	"github.com/pufferpanel/pufferpanel/v2/shared/scope"
	"github.com/spf13/viper"
	"io"
	path2 "path"
	"reflect"
	"strings"
)

func listenOnSocket(conn *websocket.Conn, server *programs.Program, scopes []scope.Scope) {
	defer func() {
		if err := recover(); err != nil {
			logging.Error("Error with websocket connection for server %s: %s", server.Id(), err)
		}
	}()

	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			logging.Exception("error on reading from websocket", err)
			return
		}
		if msgType != websocket.TextMessage {
			continue
		}
		mapping := make(map[string]interface{})

		err = json.Unmarshal(data, &mapping)
		if err != nil {
			logging.Exception("error on decoding websocket message", err)
			continue
		}

		messageType := mapping["type"]
		if message, ok := messageType.(string); ok {
			switch strings.ToLower(message) {
			case "stat":
				{
					if shared.ContainsScope(scopes, scope.ServersStat) {
						results, err := server.GetEnvironment().GetStats()
						msg := messages.StatMessage{}
						if err != nil {
							msg.Cpu = 0
							msg.Memory = 0
						} else {
							msg.Cpu = results.Cpu
							msg.Memory = results.Memory
						}
						_ = messages.Write(conn, msg)
					}
				}
			case "start":
				{
					if shared.ContainsScope(scopes, scope.ServersStart) {
						_ = server.Start()
					}
					break
				}
			case "stop":
				{
					if shared.ContainsScope(scopes, scope.ServersStop) {
						_ = server.Stop()
					}
				}
			case "install":
				{
					if shared.ContainsScope(scopes, scope.ServersInstall) {
						_ = server.Install()
					}
				}
			case "kill":
				{
					if shared.ContainsScope(scopes, scope.ServersStop) {
						_ = server.Kill()
					}
				}
			case "reload":
				{
					if shared.ContainsScope(scopes, scope.ServersEditAdmin) {
						_ = programs.Reload(server.Id())
					}
				}
			case "ping":
				{
					_ = messages.Write(conn, messages.PongMessage{})
				}
			case "console":
				{
					cmd, ok := mapping["command"].(string)
					if ok {
						if run, _ := server.IsRunning(); run {
							_ = server.GetEnvironment().ExecuteInMainProcess(cmd)
						}
					}
				}
			case "file":
				{
					if !shared.ContainsScope(scopes, scope.ServersFilesGet) {
						break
					}

					action, ok := mapping["action"].(string)
					if !ok {
						break
					}
					path, ok := mapping["path"].(string)
					if !ok {
						break
					}

					switch strings.ToLower(action) {
					case "get":
						{
							editMode, ok := mapping["edit"].(bool)
							handleGetFile(conn, server, path, ok && editMode)
						}
						break
					case "delete":
						{
							if !shared.ContainsScope(scopes, scope.ServersFilesPut) {
								break
							}

							err := server.DeleteItem(path)
							if err != nil {
								_ = messages.Write(conn, messages.FileListMessage{Error: err.Error()})
							} else {
								//now get the root
								handleGetFile(conn, server, path2.Dir(path), false)
							}
						}
						break
					case "create":
						{
							if !shared.ContainsScope(scopes, scope.ServersFilesPut) {
								break
							}

							err := server.CreateFolder(path)

							if err != nil {
								_ = messages.Write(conn, messages.FileListMessage{Error: err.Error()})
							} else {
								handleGetFile(conn, server, path, false)
							}
						}
						break
					}
				}
			default:
				_ = conn.WriteJSON(map[string]string{"error": "unknown command"})
			}
		} else {
			logging.Error("message type is not a string, but was %s", reflect.TypeOf(messageType))
		}
	}
}

func handleGetFile(conn *websocket.Conn, server *programs.Program, path string, editMode bool) {
	data, err := server.GetItem(path)
	if err != nil {
		_ = messages.Write(conn, messages.FileListMessage{Error: err.Error()})
		return
	}

	defer shared.Close(data.Contents)

	if data.FileList != nil {
		_ = messages.Write(conn, messages.FileListMessage{FileList: data.FileList, CurrentPath: path})
	} else if data.Contents != nil {
		//if the file is small enough, we'll send it over the websocket
		if editMode && data.ContentLength < viper.GetInt64("data.maxWSDownloadSize") {
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, data.Contents)
			_ = messages.Write(conn, messages.FileListMessage{Contents: buf.Bytes(), Filename: data.Name})
		} else {
			_ = messages.Write(conn, messages.FileListMessage{Url: path, Filename: data.Name})
		}
	}
}
