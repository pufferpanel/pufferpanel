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

package socket

import (
	"github.com/gorilla/websocket"
	"github.com/pufferpanel/pufferpanel/v2/daemon/messages"
)

func Write(c *websocket.Conn, msg messages.Message) error {
	return c.WriteJSON(messages.Transmission{Type: msg.Key(), Message: msg})
}
