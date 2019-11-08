/*
 Copyright 2018 Padduck, LLC

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

package messages

import (
	"github.com/gorilla/websocket"
)

type Message interface {
	Key() string
}

type Transmission struct {
	Message Message `json:"data"`
	Type    string  `json:"type"`
}

func Write(c *websocket.Conn, msg Message) error {
	return c.WriteJSON(Transmission{Type: msg.Key(), Message: msg})
}

type StatMessage struct {
	Memory float64 `json:"memory"`
	Cpu    float64 `json:"cpu"`
}

type ConsoleMessage struct {
	Logs []string `json:"logs"`
}

type PingMessage struct {
}

type PongMessage struct {
}

type FileListMessage struct {
	CurrentPath string     `json:"path"`
	Error       string     `json:"error,omitempty"`
	Url         string     `json:"url,omitempty"`
	FileList    []FileDesc `json:"files,omitempty"`
	Contents    []byte     `json:"contents,omitempty"`
	Filename    string     `json:"name,omitempty"`
}

func (m StatMessage) Key() string {
	return "stat"
}

func (m ConsoleMessage) Key() string {
	return "console"
}

func (m PingMessage) Key() string {
	return "ping"
}

func (m PongMessage) Key() string {
	return "pong"
}

func (m FileListMessage) Key() string {
	return "file"
}
