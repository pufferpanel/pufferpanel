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
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/v2/models"
)

type OAuth2 struct {
	DB *gorm.DB
}

func (o *OAuth2) Get(clientId string) (*models.Client, error) {
	client := &models.Client{
		ClientId: clientId,
	}
	err := o.DB.Where(client).Find(client).Error
	return client, err
}

func (o *OAuth2) GetForUser(userId uint) ([]*models.Client, error) {
	clients := &models.Clients{}

	client := &models.Client{
		UserId: userId,
	}

	err := o.DB.Where(client).Find(clients).Error
	return *clients, err
}

func (o *OAuth2) Update(client *models.Client) error {
	o.DB.Save(client)
	return nil
}

func (o *OAuth2) Delete(client *models.Client) error {
	o.DB.Where(client).Delete(client)
	return nil
}
