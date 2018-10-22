package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             uint   `gorm:"PRIMARY_KEY;AUTO_INCREMENT",json:"-"`
	Username       string `gorm:"UNIQUE_INDEX;NOT NULL",json:"-"`
	Email          string `gorm:"UNIQUE_INDEX;NOT NULL",json:"-"`
	HashedPassword string `gorm:"column:password;NOT NULL",json:"-"`

	//CreatedAt time.Time `json:"-"`
	//UpdatedAt time.Time `json:"-"`
}

type Users []*User

func (u *User) SetPassword(pw string) error {
	res, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)

	if err != nil {
		u.HashedPassword = string(res)
	}

	return err
}
