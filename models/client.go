package models

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
)

type Client struct {
	ID                 uint   `gorm:"PRIMARY_KEY,AUTO_INCREMEMT" json:"-"`
	ClientId           string `gorm:"NOT NULL" json:"client_id"`
	HashedClientSecret string `gorm:"column:client_secret;NOT NULL" json:"-"`

	UserId uint  `gorm:"NOT NULL" json:"-"`
	User   *User `json:"-"`

	ServerId string  `gorm:"NOT NULL" json:"-"`
	Server   *Server `json:"-"`

	Name        string `gorm:"column:name;NOT NULL;size:100;default\"\"" json:"name"`
	Description string `gorm:"column:description;NOT NULL;size:4000;default:\"\"" json:"description"`
} //@name Client

type CreatedClient struct {
	ClientId     string `json:"id"`
	ClientSecret string `json:"secret"`
} //@name CreatedClient

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
