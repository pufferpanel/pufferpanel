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
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/models"
	"io"
	"net/http"
)

var nodeClient = http.Client{}

type NodeService interface {
	GetAll() (*models.Nodes, error)

	Get(id uint) (*models.Node, bool, error)

	Update(model *models.Node) error

	Delete(id uint) error

	Create(node *models.Node) error

	CallNode(node *models.Node, method string, path string, body io.ReadCloser, headers http.Header) (*http.Response, error)
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
	_, err := createNodeURL(node, path)
	if err != nil {
		return nil, err
	}

	request := &http.Request{
		Method: method,
		URL: nil,
		Header: headers,
	}

	if method != "GET" && body != nil {
		request.Body = body
	}

	response, err := nodeClient.Do(request)
	return response, err
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

	protocol := "http"
	if ssl {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s:%d/%s", protocol, node.PrivateHost, node.PrivatePort, path), nil
}