package models

import (
	"github.com/pufferpanel/pufferpanel/shared"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

type User struct {
	ID             uint   `json:"-"`
	Username       string `gorm:"UNIQUE_INDEX;NOT NULL" json:"-" validate:"required,printascii"`
	Email          string `gorm:"UNIQUE_INDEX;NOT NULL" json:"-" validate:"required,email"`
	HashedPassword string `gorm:"column:password;NOT NULL" json:"-" validate:"required"`

	//CreatedAt time.Time `json:"-"`
	//UpdatedAt time.Time `json:"-"`
}

type Users []*User

func (u *User) SetPassword(pw string) error {
	res, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)

	if err == nil {
		u.HashedPassword = string(res)
	}

	return err
}

func (u *User) IsValid() (err error) {
	err = validator.New().Struct(u)
	if err != nil {
		err = shared.GenerateValidationMessage(err)
	}
	return
}

func (u *User) BeforeSave() (err error) {
	err = u.IsValid()
	return
}
