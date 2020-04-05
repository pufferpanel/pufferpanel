package pufferpanel

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"sync"
)

func Create(ws *websocket.Conn) *Socket {
	return &Socket{conn: ws}
}

type Socket struct {
	conn *websocket.Conn
	locker sync.Mutex
}

func (s *Socket) WriteMessage(data []byte) error {
	s.locker.Lock()
	defer s.locker.Unlock()
	return s.conn.WriteMessage(websocket.TextMessage, data)
}

func (s *Socket) WriteJSON(data interface{}) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return s.WriteMessage(d)
}

func (s *Socket) Close() error {
	return s.conn.Close()
}

func (s *Socket) ReadMessage() (messageType int, p []byte, err error) {
	return s.conn.ReadMessage()
}