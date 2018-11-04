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
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/models"
)

type NodeService interface {
	GetAll() (*models.Nodes, error)

	Get(id uint) (*models.Node, bool, error)

	Update(model *models.Node) error

	Delete(id uint) error

	Create(node *models.Node) error
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
