package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	uuid "github.com/gofrs/uuid/v5"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"gorm.io/gorm"
	"strings"
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

	sessionToken := token.String()

	res, err := hashToken(sessionToken)
	if err != nil {
		return "", err
	}

	session := &models.Session{
		Token:          res,
		ExpirationTime: time.Now().Add(time.Hour),
		UserId:         &user.ID,
	}

	err = ss.DB.Create(session).Error
	return sessionToken, err
}

func (ss *Session) CreateForClient(node *models.Client) (string, error) {
	token, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	sessionToken := token.String()

	res, err := hashToken(sessionToken)
	if err != nil {
		return "", err
	}

	session := &models.Session{
		Token:          res,
		ExpirationTime: time.Now().Add(time.Hour),
		ClientId:       &node.ID,
	}

	err = ss.DB.Create(session).Error
	return sessionToken, err
}

func (ss *Session) Validate(token string) (*models.Session, error) {
	hashed, err := hashToken(token)
	if err != nil {
		return nil, err
	}

	session := &models.Session{Token: hashed}
	query := ss.DB.Preload("Client").Preload("User").Preload("Server")
	query = query.Where("expiration_time > ?", time.Now())
	query = query.Where("user_id IS NOT NULL OR client_id IS NOT NULL")
	query = query.Where(session)

	err = query.Find(session).Error

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
	hashed, err := hashToken(token)
	if err != nil {
		return err
	}

	session := &models.Session{Token: hashed}
	err = ss.DB.Where(session).Delete(session).Error
	if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

func hashToken(source string) (result string, err error) {
	h := sha256.New()
	_, err = h.Write([]byte(source))
	if err != nil {
		return
	}
	bs := h.Sum(nil)
	builder := &strings.Builder{}
	_, err = hex.NewEncoder(builder).Write(bs)
	result = builder.String()
	return
}
