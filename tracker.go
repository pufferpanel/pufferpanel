package pufferpanel

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Tracker struct {
	sockets    map[*Socket]bool
	broadcast  chan []byte
	register   chan *Socket
	unregister chan *Socket
}

func CreateTracker() *Tracker {
	tracker := &Tracker{
		sockets:    make(map[*Socket]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Socket),
		unregister: make(chan *Socket),
	}

	go tracker.pump()

	return tracker
}

func (tracker *Tracker) pump() {
	for {
		select {
		case conn := <-tracker.register:
			tracker.sockets[conn] = true
		case conn := <-tracker.unregister:
			if _, ok := tracker.sockets[conn]; ok {
				delete(tracker.sockets, conn)
				conn.Close()
			}
		case message := <-tracker.broadcast:
			for socket := range tracker.sockets {
				select {
				case socket.send <- message:
					// message sent to client's channel
				default:
					// client too slow, disconnect
					close(socket.send)
					delete(tracker.sockets, socket)
					socket.Close()
				}
			}
		}
	}
}

func (socket *Socket) WritePump() {
	for msg := range socket.send {
		err := socket.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
	socket.Close()
}

func (tracker *Tracker) Register(conn *Socket) {
	tracker.register <- conn
}

func (ws *Tracker) WriteMessage(msg Transmission) error {
	d, err := json.Marshal(&msg)
	if err != nil {
		return err
	}
	ws.broadcast <- d
	return nil
}

func (ws *Tracker) Write(source []byte) (n int, e error) {
	packet := ServerLogs{Logs: source}
	e = ws.WriteMessage(Transmission{
		Message: packet,
		Type:    MessageTypeLog,
	})
	n = len(source)
	return
}

func Create(ws *websocket.Conn) *Socket {
	return &Socket{
		conn: ws,
		send: make(chan []byte, 256),
	}
}

type Socket struct {
	conn *websocket.Conn
	send chan []byte
}

func (s *Socket) WriteMessage(msg Transmission) error {
	return s.WriteJSON(&msg)
}

func (s *Socket) Write(data []byte) (int, error) {
	s.send <- data
	return len(data), nil
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
