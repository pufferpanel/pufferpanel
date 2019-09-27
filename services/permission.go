package services

import (
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/v2/models"
)

type Permission struct {
	DB *gorm.DB
}

func (p *Permission) GetForUser(id uint) ([]*models.Permissions, error) {
	allPerms := models.MultiplePermissions{}
	permissions := &models.Permissions{
		UserId: &id,
	}

	err := p.DB.Preload("User").Where(permissions).Find(allPerms).Error

	return allPerms, err
}

func (p *Permission) GetForUserAndServer(userId uint, serverId *string) (*models.Permissions, error) {
	permissions := &models.Permissions{
		UserId:           &userId,
		ServerIdentifier: serverId,
	}

	err := p.DB.Preload("User").Preload("Server").Where(permissions).Find(permissions).Error

	return permissions, err
}

func (p *Permission) GetForClient(id uint) ([]*models.Permissions, error) {
	allPerms := models.MultiplePermissions{}

	permissions := &models.Permissions{
		ClientId: &id,
	}

	err := p.DB.Preload("ClientId").Preload("User").Preload("Server").Where(permissions).Find(allPerms).Error

	return allPerms, err
}

func (p *Permission) GetForClientAndServer(id uint, serverId *string) (*models.Permissions, error) {
	permissions := &models.Permissions{
		ClientId:         &id,
		ServerIdentifier: serverId,
	}

	err := p.DB.Preload("User").Preload("Server").Where(permissions).Find(permissions).Error

	return permissions, err
}

