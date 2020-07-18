/*
 Copyright 2020 Padduck, LLC
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
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/v2/models"
	uuid2 "github.com/satori/go.uuid"
	"strings"
)

type Server struct {
	DB *gorm.DB
}

type ServerSearch struct {
	Username string
	NodeId   uint
	NodeName string
	Name     string
	PageSize uint
	Page     uint
}

func (ss *Server) Search(searchCriteria ServerSearch) (records *models.Servers, total uint, err error) {
	records = &models.Servers{}

	query := ss.DB

	if searchCriteria.NodeId != 0 {
		query = query.Where(&models.Server{NodeID: searchCriteria.NodeId})
	} else if searchCriteria.NodeName != "" {
		query = query.Joins("JOIN nodes n ON servers.node_id = n.id AND n.name = ?", searchCriteria.NodeName)
	}

	if searchCriteria.Username != "" {
		query = query.Joins("JOIN permissions p ON servers.identifier = p.server_identifier AND p.view_server = 1")
		query = query.Joins("JOIN users ON p.user_id = users.id")
		query = query.Where("users.username = ?", searchCriteria.Username)
	}

	nameFilter := strings.Replace(searchCriteria.Name, "*", "%", -1)

	if nameFilter != "" && nameFilter != "%" {
		query = query.Where("name LIKE ?", nameFilter)
	}

	err = query.Model(&records).Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Node").Offset((searchCriteria.Page - 1) * searchCriteria.PageSize).Limit(searchCriteria.PageSize).Order("servers.name").Find(records).Error

	return
}

func (ss *Server) Get(id string) (*models.Server, error) {
	if id == "" {
		return nil, gorm.ErrRecordNotFound
	}
	model := &models.Server{
		Identifier: id,
	}

	err := ss.DB.Where(model).Preload("Node").First(model).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (ss *Server) Update(model *models.Server) error {
	res := ss.DB.Save(model)
	return res.Error
}

func (ss *Server) Delete(id string) error {
	//if we are already in a transaction, use the existing transaction
	inTrans := false
	var trans *gorm.DB
	if _, ok := ss.DB.CommonDB().(*sql.Tx); ok {
		trans = ss.DB
		inTrans = true
	} else {
		trans = ss.DB.Begin()
		defer trans.RollbackUnlessCommitted()
	}

	model := &models.Server{
		Identifier: id,
	}

	err := trans.Delete(models.Permissions{}, "server_identifier = ?", id).Error
	if err != nil {
		return err
	}

	err = trans.Delete(models.Client{}, "server_id = ?", id).Error
	if err != nil {
		return err
	}

	err = ss.DB.Delete(model).Error
	if err != nil {
		return err
	}

	if inTrans {
		return nil
	} else {
		return trans.Commit().Error
	}
}

func (ss *Server) Create(model *models.Server) (err error) {
	if model.Identifier == "" {
		uuid := uuid2.NewV4()
		generatedId := strings.ToUpper(uuid.String())[0:8]
		model.Identifier = generatedId
	}

	res := ss.DB.Create(model)
	if res.Error != nil {
		err = res.Error
		return
	}
	return
}
