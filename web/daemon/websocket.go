/*
 Copyright 2020 Padduck, LLC
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at
  	http://www.apache.org/licenses/LICENSE-2.0
  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package daemon

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/messages"
	"github.com/pufferpanel/pufferpanel/v3/servers"
	"reflect"
	"runtime/debug"
	"strings"
)

func listenOnSocket(conn *pufferpanel.Socket, server *servers.Server, scopes []*pufferpanel.Scope) {
	defer func() {
		if err := recover(); err != nil {
			logging.Error.Printf("Error with websocket connection for server %s: %s\n%s", server.Id(), err, debug.Stack())
		}
	}()

	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			logging.Error.Printf("error on reading from websocket: %s", err)
			return
		}
		if msgType != websocket.TextMessage {
			continue
		}
		mapping := make(map[string]interface{})

		err = json.Unmarshal(data, &mapping)
		if err != nil {
			logging.Error.Printf("error on decoding websocket message: %s", err)
			continue
		}

		messageType := mapping["type"]
		if message, ok := messageType.(string); ok {
			switch strings.ToLower(message) {
			case "replay":
				{
					if pufferpanel.ContainsScope(scopes, pufferpanel.ScopeServerConsole) {
						console, _ := server.GetEnvironment().GetConsole()
						_ = conn.WriteMessage(messages.Console{Logs: console})
					}
				}

			default:
				_ = conn.WriteJSON(map[string]string{"error": "unknown command"})
			}
		} else {
			logging.Error.Printf("message type is not a string, but was %s", reflect.TypeOf(messageType))
		}
	}
}
