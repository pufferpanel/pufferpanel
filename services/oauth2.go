package services

import (
	"github.com/pufferpanel/pufferpanel/v3/models"
	"gorm.io/gorm"
)

type OAuth2 struct {
	DB *gorm.DB
}

// Get Gets a specific OAuth client, including the scopes it holds
func (s *OAuth2) Get(clientId string) (*models.Client, error) {
	client := &models.Client{
		ClientId: clientId,
	}
	err := s.DB.Where(client).First(client).Error
	return client, err
}

// GetForUser Gets all clients for a user
func (s *OAuth2) GetForUser(userId uint) ([]*models.Client, error) {
	client := &models.Client{
		UserId: userId,
	}
	var clients []*models.Client
	err := s.DB.Where(client).Find(&clients).Error
	return clients, err
}

func (s *OAuth2) Create(request *models.Client) error {
	return s.DB.Create(request).Error
}

func (s *OAuth2) Update(request *models.Client) error {
	return s.DB.Save(request).Error
}

func (s *OAuth2) Delete(clientId string) error {
	client := &models.Client{
		ClientId: clientId,
	}
	return s.DB.Model(client).Delete(client, client).Error
}
