package services

import (
	uuid "github.com/gofrs/uuid/v5"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"gorm.io/gorm"
	"time"
)

type Session struct {
	DB *gorm.DB
}

func (ss *Session) CreateForUser(user *models.User) (string, error) {
	token, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	session := &models.Session{
		Token:          token.String(),
		ExpirationTime: time.Now().Add(time.Hour),
		UserId:         &user.ID,
	}

	err = ss.DB.Create(session).Error
	return token.String(), err
}

func (ss *Session) CreateForClient(node *models.Client) (string, error) {
	token, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	session := &models.Session{
		Token:          token.String(),
		ExpirationTime: time.Now().Add(time.Hour),
		ClientId:       &node.ID,
	}

	err = ss.DB.Create(session).Error
	return token.String(), err
}

func (ss *Session) Validate(token string) (*models.Session, error) {
	session := &models.Session{Token: token}
	err := ss.DB.Preload("Client").Preload("User").Preload("Server").Where(session).Find(session).Error

	//validate this session is for a client or user
	if err == nil && session.ClientId == nil && session.UserId == nil {
		err = gorm.ErrRecordNotFound
	}

	if session.ExpirationTime.Before(time.Now()) {
		return nil, gorm.ErrRecordNotFound
	}

	return session, err
}

func (ss *Session) ValidateNode(token string) (*models.Node, error) {
	if models.LocalNode != nil && models.LocalNode.Secret == token {
		return models.LocalNode, nil
	}

	node := &models.Node{Secret: token}
	err := ss.DB.Where(node).First(node).Error
	return node, err
}

func (ss *Session) Expire(token string) error {
	session := &models.Session{Token: token}
	return ss.DB.Delete(session).Error
}
