package models

import (
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"golang.org/x/crypto/blake2b"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
)

type RecoveryCode struct {
	ID uint `gorm:"column:id;primaryKey;autoIncrement" json:"-"`

	UserId uint `gorm:"column:user_id;not null;index" json:"-"`
	User   User `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

	CodeHash string `gorm:"column:code;not null" json:"-"`
}

func (rc *RecoveryCode) SetCode(code string) error {
	hash, err := blake2b.New256(nil)
	if err != nil {
		return err
	}

	_, err = hash.Write([]byte(code))
	if err != nil {
		return err
	}
	
	rc.CodeHash = fmt.Sprintf("%x", hash.Sum(nil))
	return nil
}

func (rc *RecoveryCode) IsValid() (err error) {
	err = validator.New().Struct(rc)
	if err != nil {
		err = pufferpanel.GenerateValidationMessage(err)
	}
	return
}

func (rc *RecoveryCode) BeforeSave(*gorm.DB) error {
	return rc.IsValid()
}
