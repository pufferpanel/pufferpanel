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
	"github.com/pufferpanel/pufferpanel/v2/models"
	"gorm.io/gorm"
)

type UserSettings struct {
	DB *gorm.DB
}

func (uss *UserSettings) GetAllForUser(userId uint) (models.UserSettingsView, error) {
	records := &models.UserSettings{}

	query := uss.DB

	query = query.Where(&models.UserSetting{UserID: userId})

	err := query.Model(&records).Error
	if err != nil {
		return nil, err
	}

	err = query.Find(records).Error
	if err != nil {
		return nil, err
	}

	return models.FromUserSettings(records), nil
}

func (uss *UserSettings) Update(model *models.UserSetting) error {
	search := &models.UserSetting{
		Key:    model.Key,
		UserID: model.UserID,
	}

	err := uss.DB.Where(search).First(search).Error

	if err != nil && gorm.ErrRecordNotFound != err {
		return err
	}

	if err != nil {
		err = uss.DB.Create(model).Error
	} else {
		err = uss.DB.Save(model).Error
	}

	return err
}
