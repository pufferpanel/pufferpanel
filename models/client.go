package models

import (
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/v2"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

type Client struct {
	gorm.Model

	ClientId           string
	HashedClientSecret string

	UserId uint
	User   User
}

func (c *Client) SetClientSecret(secret string) error {
	res, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)

	if err == nil {
		c.HashedClientSecret = string(res)
	}

	return err
}

func (c *Client) IsValid() (err error) {
	err = validator.New().Struct(c)

	if err != nil {
		err = pufferpanel.GenerateValidationMessage(err)
	}

	return
}

func (c *Client) BeforeSave() (err error) {
	err = c.IsValid()
	return
}
