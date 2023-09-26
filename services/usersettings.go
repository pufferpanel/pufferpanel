package services

import (
	"github.com/pufferpanel/pufferpanel/v3/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserSettings struct {
	DB *gorm.DB
}

func (uss *UserSettings) GetAllForUser(userId uint) (models.UserSettingsView, error) {
	var records []*models.UserSetting

	query := uss.DB

	query = query.Where(&models.UserSetting{UserID: userId})

	err := query.Model(&records).Error
	if err != nil {
		return nil, err
	}

	err = query.Find(&records).Error
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
		err = uss.DB.Omit(clause.Associations).Create(model).Error
	} else {
		err = uss.DB.Omit(clause.Associations).Save(model).Error
	}

	return err
}
