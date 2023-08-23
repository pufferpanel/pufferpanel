/*
 Copyright 2023 PufferPanel

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

package pufferpanel

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/messages"
	"io"
	"sync"
)

type Tracker struct {
	sockets []*Socket
	locker  sync.Mutex
}

func CreateTracker() *Tracker {
	return &Tracker{sockets: make([]*Socket, 0)}
}

func (ws *Tracker) Register(conn *Socket) {
	ws.locker.Lock()
	defer ws.locker.Unlock()
	ws.sockets = append(ws.sockets, conn)
}

func (ws *Tracker) WriteMessage(msg messages.Message) error {
	d, err := json.Marshal(&messages.Transmission{Message: msg, Type: msg.Key()})
	if err != nil {
		return err
	}
	ws.locker.Lock()
	defer ws.locker.Unlock()

	for i := 0; i < len(ws.sockets); i++ {
		go func(conn *Socket, data []byte) {
			_, err := conn.Write(data)
			if err != nil {
				logging.Debug.Printf("websocket encountered error, dropping (%s)", err.Error())
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

func Create(ws *websocket.Conn) *Socket {
	return &Socket{conn: ws}
}

type Socket struct {
	conn   *websocket.Conn
	locker sync.Mutex
	io.WriteCloser
}

func (s *Socket) WriteMessage(msg messages.Message) error {
	return s.WriteJSON(messages.Transmission{Type: msg.Key(), Message: msg})
}

func (s *Socket) Write(data []byte) (int, error) {
	s.locker.Lock()
	defer s.locker.Unlock()
	return len(data), s.conn.WriteMessage(websocket.TextMessage, data)
}

func (s *Socket) WriteJSON(data interface{}) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = s.Write(d)
	return err
}

func (s *Socket) Close() error {
	return s.conn.Close()
}

func (s *Socket) ReadMessage() (messageType int, p []byte, err error) {
	return s.conn.ReadMessage()
}
