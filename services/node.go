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

package services

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

var nodeClient = http.Client{}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Node struct {
	DB *gorm.DB
}

func (ns *Node) GetAll() (*models.Nodes, error) {
	nodes := &models.Nodes{}

	res := ns.DB.Find(nodes)

	return nodes, res.Error
}

func (ns *Node) Get(id uint) (*models.Node, error) {
	model := &models.Node{}

	res := ns.DB.First(model, id)

	return model, res.Error
}

func (ns *Node) Update(model *models.Node) error {
	res := ns.DB.Save(model)
	return res.Error
}

func (ns *Node) Delete(id uint) error {
	model := &models.Node{
		ID: id,
	}

	var count int
	ns.DB.Model(&models.Server{}).Where("node_id = ?", model.ID).Count(&count)
	if count > 0 {
		return pufferpanel.ErrNodeHasServers
	}

	res := ns.DB.Delete(model)
	return res.Error
}

func (ns *Node) Create(node *models.Node) error {
	res := ns.DB.Create(node)
	return res.Error
}

func (ns *Node) CallNode(node *models.Node, method string, path string, body io.ReadCloser, headers http.Header) (*http.Response, error) {
	var fullUrl string
	var err error

	if node.IsLocal() {
		fullUrl = "http://localhost" + path
	} else {
		fullUrl, err = createNodeURL(node, path)
		if err != nil {
			return nil, err
		}
	}

	addr, err := url.Parse(fullUrl)
	if err != nil {
		return nil, err
	}

	request := &http.Request{
		Method: method,
		URL:    addr,
		Header: headers,
	}

	if method != "GET" && body != nil {
		request.Body = body
	}

	if node.IsLocal() {
		w := &httptest.ResponseRecorder{}
		pufferpanel.Engine.ServeHTTP(w, request)
		return w.Result(), err
	}

	response, err := nodeClient.Do(request)
	return response, err
}

func (ns *Node) OpenSocket(node *models.Node, path string, writer http.ResponseWriter, request *http.Request) error {
	ssl, err := doesDaemonUseSSL(node)
	if err != nil {
		return err
	}

	scheme := "ws"
	if ssl {
		scheme = "wss"
	}
	addr := fmt.Sprintf("%s:%d", node.PrivateHost, node.PrivatePort)

	u := url.URL{Scheme: scheme, Host: addr, Path: path}
	logging.Debug().Printf("Proxying connection to %s", u.String())

	header := http.Header{}
	header.Set("Authorization", request.Header.Get("Authorization"))

	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		return err
	}

	conn, err := wsupgrader.Upgrade(writer, request, nil)
	if err != nil {
		return err
	}

	go func(daemon *websocket.Conn, client *websocket.Conn) {
		defer func() {
			_ = daemon.Close()
			_ = client.Close()
		}()

		ch := make(chan error)
		go proxyRead(daemon, client, ch)
		go proxyRead(client, daemon, ch)

		err := <-ch

		if err != nil {
			logging.Error().Printf("Error proxying socket: %s", err)
		}
	}(c, conn)

	return nil
}

func doesDaemonUseSSL(node *models.Node) (bool, error) {
	if node.IsLocal() {
		return false, nil
	}

	path := fmt.Sprintf("://%s:%d/daemon", node.PrivateHost, node.PrivatePort)

	_, err := http.Get("https" + path)

	if err != nil {
		_, err = http.Get("http" + path)
		return false, err
	}

	return true, nil
}

func createNodeURL(node *models.Node, path string) (string, error) {
	ssl, err := doesDaemonUseSSL(node)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}

	if strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}

	protocol := "http"
	if ssl {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s:%d/%s", protocol, node.PrivateHost, node.PrivatePort, path), nil
}

func proxyRead(source, dest *websocket.Conn, ch chan error) {
	for {
		messageType, data, err := source.ReadMessage()

		if err != nil {
			ch <- err
			return
		}
		err = dest.WriteMessage(messageType, data)
		if err != nil {
			ch <- err
			return
		}
	}
}
