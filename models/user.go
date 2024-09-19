package models

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID             uint   `gorm:"column:id;primaryKey;autoIncrement" json:"-"`
	Username       string `gorm:"column:username;not null;size:100;uniqueIndex;unique" json:"-" validate:"required,printascii,max=100,min=5"`
	Email          string `gorm:"column:email;not null;size:255;uniqueIndex;unique" json:"-" validate:"required,email,max=255"`
	HashedPassword string `gorm:"column:password;NOT NULL;size:200" json:"-" validate:"required,max=200"`
	OtpSecret      string `gorm:"column:otp_secret;size:32" json:"-"`
	OtpActive      bool   `gorm:"column:otp_active;not null;DEFAULT:0" json:"-"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

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
		err = pufferpanel.GenerateValidationMessage(err)
	}
	return
}

func (u *User) BeforeSave(*gorm.DB) (err error) {
	err = u.IsValid()
	return
}
