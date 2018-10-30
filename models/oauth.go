package models

import (
	"github.com/pufferpanel/pufferpanel/shared"
	"gopkg.in/go-playground/validator.v9"
	"strings"
	"time"
)

type OauthClient struct {
	ID uint `json:"-"`

	ClientId     string
	ClientSecret string

	UserID uint `json:"-"`
	User   User `gorm:"association_autoupdate:false" json:"-"`

	NodeID uint `json:"-"`
	Node   Node `gorm:"association_autoupdate:false" json:"-"`

	Scopes   []string `gorm:"-" json:"-" validate:"unique"`
	dbScopes string   `gorm:"NOT NULL" json:"-" validate:"required"`

	AccessTokens []OauthAccessToken
}

type OauthAccessToken struct {
	ID            uint `json:"-"`
	OauthClientID uint `gorm:"NOT NULL" json:"-" validate:"required"`

	Token      string    `gorm:"NOT NULL;UNIQUE_INDEX" json:"-" validate:"required"`
	ExpireDate time.Time `gorm:"NOT NULL" json:"-" validate:"required"`
}

func (o *OauthClient) BeforeSave() (err error) {
	o.dbScopes = strings.Join(o.Scopes, ",")
	return
}

func (o *OauthClient) AfterFind() (err error) {
	o.Scopes = strings.Split(o.dbScopes, ",")
	return
}

func (m *OauthClient) IsValid() (err error) {
	err = validator.New().Struct(m)
	if err != nil {
		err = shared.GenerateValidationMessage(err)
	}
	return
}

func (m *OauthAccessToken) IsValid() (err error) {
	err = validator.New().Struct(m)
	if err != nil {
		err = shared.GenerateValidationMessage(err)
	}
	return
}
