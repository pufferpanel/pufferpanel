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

package socket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pufferpanel/pufferpanel/v2/daemon/messages"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"sync"
)

type Tracker struct {
	sockets []*connection
	locker  sync.Mutex
}

type connection struct {
	socket *websocket.Conn
	lock sync.Mutex
}

func CreateTracker() *Tracker {
	return &Tracker{sockets: make([]*connection, 0)}
}

func (ws *Tracker) Register(conn *websocket.Conn) {
	ws.locker.Lock()
	defer ws.locker.Unlock()
	ws.sockets = append(ws.sockets, &connection{socket: conn})
}

func (ws *Tracker) WriteMessage(msg messages.Message) error {
	d, err := json.Marshal(&messages.Transmission{Message: msg, Type: msg.Key()})
	if err != nil {
		return err
	}
	ws.locker.Lock()
	defer ws.locker.Unlock()

	for i := 0; i < len(ws.sockets); i++ {
		go func(conn *connection, data []byte) {
			conn.lock.Lock()
			defer conn.lock.Unlock()
			err := conn.socket.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				logging.Info().Printf("websocket encountered error, dropping (%s)", err.Error())
				ws.locker.Lock()
				defer ws.locker.Unlock()
				for i, k := range ws.sockets {
					if k == conn {
						ws.sockets[i] = ws.sockets[len(ws.sockets)-1]
						ws.sockets[len(ws.sockets)-1] = nil
						ws.sockets = ws.sockets[:len(ws.sockets)-1]
						break
					}
				}
			}
		}(ws.sockets[i], d)
	}

	return nil
}

func (ws *Tracker) Write(source []byte) (n int, e error) {
	logs := make([]string, 1)
	logs[0] = string(source)
	packet := messages.Console{Logs: logs}
	e = ws.WriteMessage(packet)
	n = len(source)
	return
}
