package server

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/messages"
	"github.com/pufferpanel/pufferpanel/v2/programs"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/spf13/viper"
	"io"
	path2 "path"
	"reflect"
	"runtime/debug"
	"strings"
)

func listenOnSocket(conn *websocket.Conn, server *programs.Program, scopes []pufferpanel.Scope) {
	defer func() {
		if err := recover(); err != nil {
			logging.Error().Printf("Error with websocket connection for server %s: %s\n%s", server.Id(), err, debug.Stack())
		}
	}()

	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			logging.Error().Printf("error on reading from websocket: %s", err)
			return
		}
		if msgType != websocket.TextMessage {
			continue
		}
		mapping := make(map[string]interface{})

		err = json.Unmarshal(data, &mapping)
		if err != nil {
			logging.Error().Printf("error on decoding websocket message: %s", err)
			continue
		}

		messageType := mapping["type"]
		if message, ok := messageType.(string); ok {
			switch strings.ToLower(message) {
			case "stat":
				{
					if pufferpanel.ContainsScope(scopes, pufferpanel.ScopeServersStat) {
						results, err := server.GetEnvironment().GetStats()
						msg := messages.Stat{}
						if err != nil {
							msg.Cpu = 0
							msg.Memory = 0
						} else {
							msg.Cpu = results.Cpu
							msg.Memory = results.Memory
						}
						_ = pufferpanel.Write(conn, msg)
					}
				}
			case "status": {
				running, err := server.IsRunning()
				if err != nil {
					running = false
				}
				msg := messages.Status{Running:running}
				_ = pufferpanel.Write(conn, msg)
			}
			case "start":
				{
					if pufferpanel.ContainsScope(scopes, pufferpanel.ScopeServersStart) {
						_ = server.Start()
					}
					break
				}
			case "stop":
				{
					if pufferpanel.ContainsScope(scopes, pufferpanel.ScopeServersStop) {
						_ = server.Stop()
					}
				}
			case "install":
				{
					if pufferpanel.ContainsScope(scopes, pufferpanel.ScopeServersInstall) {
						_ = server.Install()
					}
				}
			case "kill":
				{
					if pufferpanel.ContainsScope(scopes, pufferpanel.ScopeServersStop) {
						_ = server.Kill()
					}
				}
			case "reload":
				{
					if pufferpanel.ContainsScope(scopes, pufferpanel.ScopeServersEditAdmin) {
						_ = programs.Reload(server.Id())
					}
				}
			case "ping":
				{
					_ = pufferpanel.Write(conn, messages.Pong{})
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
					if !pufferpanel.ContainsScope(scopes, pufferpanel.ScopeServersFilesGet) {
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
							if !pufferpanel.ContainsScope(scopes, pufferpanel.ScopeServersFilesPut) {
								break
							}

							err := server.DeleteItem(path)
							if err != nil {
								_ = pufferpanel.Write(conn, messages.FileList{Error: err.Error()})
							} else {
								//now get the root
								handleGetFile(conn, server, path2.Dir(path), false)
							}
						}
						break
					case "create":
						{
							if !pufferpanel.ContainsScope(scopes, pufferpanel.ScopeServersFilesPut) {
								break
							}

							err := server.CreateFolder(path)

							if err != nil {
								_ = pufferpanel.Write(conn, messages.FileList{Error: err.Error()})
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
			logging.Error().Printf("message type is not a string, but was %s", reflect.TypeOf(messageType))
		}
	}
}

func handleGetFile(conn *websocket.Conn, server *programs.Program, path string, editMode bool) {
	data, err := server.GetItem(path)
	if err != nil {
		_ = pufferpanel.Write(conn, messages.FileList{Error: err.Error()})
		return
	}

	defer pufferpanel.Close(data.Contents)

	if data.FileList != nil {
		_ = pufferpanel.Write(conn, messages.FileList{FileList: data.FileList, CurrentPath: path})
	} else if data.Contents != nil {
		//if the file is small enough, we'll send it over the websocket
		if editMode && data.ContentLength < viper.GetInt64("daemon.data.maxWSDownloadSize") {
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, data.Contents)
			_ = pufferpanel.Write(conn, messages.FileList{Contents: buf.Bytes(), Filename: data.Name})
		} else {
			_ = pufferpanel.Write(conn, messages.FileList{Url: path, Filename: data.Name})
		}
	}
}
