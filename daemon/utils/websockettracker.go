/*
 Copyright 2016 Padduck, LLC

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

package utils

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pufferpanel/pufferpanel/v2/daemon/messages"
	"github.com/pufferpanel/pufferpanel/v2/shared/logging"
	"sync"
)

type WebSocketManager interface {
	Register(ws *websocket.Conn)

	Write(msg []byte) (n int, e error)
}

type wsManager struct {
	sockets []websocket.Conn
	locker  sync.Mutex
}

func CreateWSManager() WebSocketManager {
	return &wsManager{sockets: make([]websocket.Conn, 0), locker: sync.Mutex{}}
}

func (ws *wsManager) Register(conn *websocket.Conn) {
	ws.sockets = append(ws.sockets, *conn)
}

func (ws *wsManager) Write(source []byte) (n int, e error) {
	ws.locker.Lock()
	logs := make([]string, 1)
	logs[0] = string(source)
	packet := messages.ConsoleMessage{Logs: logs}
	data, _ := json.Marshal(&messages.Transmission{Message: packet, Type: packet.Key()})

	for i := 0; i < len(ws.sockets); i++ {
		socket := ws.sockets[i]
		err := socket.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			logging.Debug("websocket encountered error, dropping (%s)", err.Error())
			if i+1 == len(ws.sockets) {
				ws.sockets = ws.sockets[:i]
			} else {
				ws.sockets = append(ws.sockets[:i], ws.sockets[i+1:]...)
			}
			i--
		}
	}
	ws.locker.Unlock()

	n = len(source)
	return
}
