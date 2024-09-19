package models

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
)

type Client struct {
	ID                 uint   `gorm:"column:id;primaryKey;autoIncrement" json:"-"`
	ClientId           string `gorm:"column:client_id;not null;size:100;uniqueIndex;unique" json:"client_id"`
	HashedClientSecret string `gorm:"column:client_secret;not null;size:100" json:"-"`

	ClientSecret string `gorm:"-" json:"client_secret"`

	UserId uint  `gorm:"column:user_id;not null;index" json:"-"`
	User   *User `json:"-"`

	ServerId *uint   `gorm:"column:server_id" json:"-"`
	Server   *Server `json:"-"`

	Name        string `gorm:"column:name;not null;size:100;default:''" json:"name"`
	Description string `gorm:"column:description;not null;size:4000;default:''" json:"description"`
}

func (c *Client) SetClientSecret(secret string) error {
	res, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)

	if err == nil {
		c.HashedClientSecret = string(res)
	}

	return err
}

func (c *Client) ValidateSecret(secret string) bool {
	return bcrypt.CompareHashAndPassword([]byte(c.HashedClientSecret), []byte(secret)) == nil
}

func (c *Client) IsValid() (err error) {
	err = validator.New().Struct(c)

	if err != nil {
		err = pufferpanel.GenerateValidationMessage(err)
	}

	return
}

func (c *Client) BeforeSave(*gorm.DB) error {
	return c.IsValid()
}
