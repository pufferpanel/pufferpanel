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

type Permission struct {
	DB *gorm.DB
}

func (ps *Permission) GetForUser(id uint) ([]*models.Permissions, error) {
	allPerms := &models.MultiplePermissions{}
	permissions := &models.Permissions{
		UserId: &id,
	}

	err := ps.DB.Preload("User").Preload("Server").Where(permissions).Find(&allPerms).Error

	return *allPerms, err
}

func (ps *Permission) GetForServer(serverId string) ([]*models.Permissions, error) {
	allPerms := &models.MultiplePermissions{}
	permissions := &models.Permissions{
		ServerIdentifier: &serverId,
	}

	err := ps.DB.Preload("User").Preload("Server").Where(permissions).Find(&allPerms).Error

	return *allPerms, err
}

func (ps *Permission) GetForUserAndServer(userId uint, serverId *string) (*models.Permissions, error) {
	permissions := &models.Permissions{
		UserId:           &userId,
		ServerIdentifier: serverId,
	}

	err := ps.DB.Preload("User").Preload("Server").Where(permissions).First(permissions).Error

	if err != nil && gorm.IsRecordNotFoundError(err) {
		return permissions, nil
	}

	return permissions, err
}

func (ps *Permission) GetForClient(id uint) ([]*models.Permissions, error) {
	allPerms := &models.MultiplePermissions{}

	permissions := &models.Permissions{
		ClientId: &id,
	}

	err := ps.DB.Preload("ClientId").Preload("User").Preload("Server").Where(permissions).Find(&allPerms).Error

	return *allPerms, err
}

func (ps *Permission) GetForClientAndServer(id uint, serverId *string) (*models.Permissions, error) {
	permissions := &models.Permissions{
		ClientId:         &id,
		ServerIdentifier: serverId,
	}

	err := ps.DB.Preload("User").Preload("Server").Where(permissions).FirstOrCreate(permissions).Error

	return permissions, err
}

func (ps *Permission) UpdatePermissions(perms *models.Permissions) error {
	//update oauth2 with new information
	if perms.ShouldDelete() {
		return ps.Remove(perms)
	} else {
		return ps.DB.Save(perms).Error
	}
}

func (ps *Permission) Remove(perms *models.Permissions) error {
	//update oauth2 with new information

	return ps.DB.Delete(perms).Error
}
