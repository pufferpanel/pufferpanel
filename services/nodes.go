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
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/models"
	"io"
	"log"
	"net/http"
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

type NodeService interface {
	GetAll() (*models.Nodes, error)

	Get(id uint) (*models.Node, bool, error)

	Update(model *models.Node) error

	Delete(id uint) error

	Create(node *models.Node) error

	CallNode(node *models.Node, method string, path string, body io.ReadCloser, headers http.Header) (*http.Response, error)

	OpenSocket(node *models.Node, path string, writer http.ResponseWriter, request *http.Request) error
}

type nodeService struct {
	db *gorm.DB
}

func GetNodeService() (NodeService, error) {
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	service := &nodeService{
		db: db,
	}

	return service, nil
}

func (ns *nodeService) GetAll() (*models.Nodes, error) {
	nodes := &models.Nodes{}

	res := ns.db.Find(nodes)

	return nodes, res.Error
}

func (ns *nodeService) Get(id uint) (*models.Node, bool, error) {
	model := &models.Node{}

	res := ns.db.FirstOrInit(model, id)

	return model, model.ID != 0, res.Error
}

func (ns *nodeService) Update(model *models.Node) error {
	res := ns.db.Save(model)
	return res.Error
}

func (ns *nodeService) Delete(id uint) error {
	model := &models.Node{
		ID: id,
	}

	res := ns.db.Delete(model)
	return res.Error
}

func (ns *nodeService) Create(node *models.Node) error {
	res := ns.db.Create(node)
	return res.Error
}

func (ns *nodeService) CallNode(node *models.Node, method string, path string, body io.ReadCloser, headers http.Header) (*http.Response, error) {
	fullUrl, err := createNodeURL(node, path)
	if err != nil {
		return nil, err
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

	response, err := nodeClient.Do(request)
	return response, err
}

func (ns *nodeService) OpenSocket(node *models.Node, path string, writer http.ResponseWriter, request *http.Request) error {
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
	log.Printf("connecting to %s", u.String())

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
		defer daemon.Close()
		defer client.Close()

		ch := make(chan error)
		go proxyRead(daemon, client, ch)
		go proxyRead(client, daemon, ch)

		err := <-ch

		if err != nil {
			logging.Exception("error proxying socket", err)
		}
	}(c, conn)

	return nil
}

func doesDaemonUseSSL(node *models.Node) (bool, error) {
	path := fmt.Sprintf("://%s:%d", node.PrivateHost, node.PrivatePort)

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
