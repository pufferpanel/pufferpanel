package services

import (
	"errors"
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
	return nil, errors.New("not implemented")
}

func (s *OAuth2) Create(request *models.Client) error {
	return errors.New("not implemented")
}

func (s *OAuth2) Update(request *models.Client) error {
	return errors.New("not implemented")
}

func (s *OAuth2) Delete(clientId string) error {
	return errors.New("not implemented")
}
