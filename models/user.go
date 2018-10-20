package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             uint   `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	Username       string `gorm:"UNIQUE_INDEX;NOT NULL"`
	Email          string `gorm:"UNIQUE_INDEX;NOT NULL"`
	HashedPassword string `gorm:"column:password;NOT NULL"`

	//CreatedAt time.Time
	//UpdatedAt time.Time
}

type Users []*User

func (u *User) SetPassword(pw string) error {
	res, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)

	if err != nil {
		u.HashedPassword = string(res)
	}

	return err
}
