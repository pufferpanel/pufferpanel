package services

import (
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"reflect"
)

type Permission struct {
	DB *gorm.DB
}

func (p *Permission) GetForUser(id uint) ([]*models.Permissions, error) {
	allPerms := &models.MultiplePermissions{}
	permissions := &models.Permissions{
		UserId: &id,
	}

	err := p.DB.Preload("User").Preload("Server").Where(permissions).Find(&allPerms).Error

	return *allPerms, err
}

func (p *Permission) GetForServer(serverId string) ([]*models.Permissions, error) {
	allPerms := &models.MultiplePermissions{}
	permissions := &models.Permissions{
		ServerIdentifier: &serverId,
	}

	err := p.DB.Preload("User").Preload("Server").Where(permissions).Find(&allPerms).Error

	return *allPerms, err
}

func (p *Permission) GetForUserAndServer(userId uint, serverId *string) (*models.Permissions, error) {
	permissions := &models.Permissions{
		UserId:           &userId,
		ServerIdentifier: serverId,
	}

	err := p.DB.Preload("User").Preload("Server").Where(permissions).First(permissions).Error

	if err != nil && gorm.IsRecordNotFoundError(err) {
		return permissions, nil
	}

	return permissions, err
}

func (p *Permission) GetForClient(id uint) ([]*models.Permissions, error) {
	allPerms := &models.MultiplePermissions{}

	permissions := &models.Permissions{
		ClientId: &id,
	}

	err := p.DB.Preload("ClientId").Preload("User").Preload("Server").Where(permissions).Find(&allPerms).Error

	return *allPerms, err
}

func (p *Permission) GetForClientAndServer(id uint, serverId *string) (*models.Permissions, error) {
	permissions := &models.Permissions{
		ClientId:         &id,
		ServerIdentifier: serverId,
	}

	err := p.DB.Preload("User").Preload("Server").Where(permissions).FirstOrCreate(permissions).Error

	return permissions, err
}

func (p *Permission) UpdatePermissions(perms *models.Permissions) error {
	shouldDelete := true

	t := reflect.TypeOf(p)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		_, exists := f.Tag.Lookup("oneOf")

		fv := reflect.ValueOf(p).FieldByName(f.Name)

		if exists && f.Type.Name() == "bool" && fv.Bool() {
			shouldDelete = false
			break
		}
	}

	if shouldDelete {
		return p.Remove(perms)
	} else {
		return p.DB.Save(perms).Error
	}
}

func (p *Permission) Remove(perms *models.Permissions) error {
	return p.DB.Where(perms).Delete(perms).Error
}
