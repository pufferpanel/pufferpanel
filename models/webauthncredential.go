package models

import (
	"github.com/go-webauthn/webauthn/webauthn"
	webauthnProto "github.com/go-webauthn/webauthn/protocol"
	"github.com/pufferpanel/pufferpanel/v3"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
)

type WebauthnCredential struct {
	ID   webauthnProto.URLEncodedBase64 `gorm:"column:id;not null;type:varchar(96);serializer:json;primaryKey" json:"-"`
	Name string                         `gorm:"column:name;not null" json:"-"`

	UserId uint `gorm:"column:user_id;not null;index" json:"-"`
	User   User `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

	Credential webauthn.Credential `gorm:"column:credential;not null;serializer:json" json:"-"`
}

func (wac *WebauthnCredential) IsValid() (err error) {
	err = validator.New().Struct(wac)
	if err != nil {
		err = pufferpanel.GenerateValidationMessage(err)
	}
	return
}

func (wac *WebauthnCredential) BeforeSave(*gorm.DB) error {
	return wac.IsValid()
}
