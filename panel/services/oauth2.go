package services

import (
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/v2/panel/models"
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
