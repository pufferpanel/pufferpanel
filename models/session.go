package models

import (
	"time"
)

type Session struct {
	ID             uint      `gorm:"primaryKey,autoIncrement" json:"-"`
	Token          string    `gorm:"unique;size:36;not null" json:"-"`
	ExpirationTime time.Time `gorm:"not null" json:"-"`

	UserId *uint `json:"-"`
	User   User  `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

	ClientId *uint  `json:"-"`
	Client   Client `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

	//if this set is for a server, what server
	ServerIdentifier *string `json:"-"`
	Server           Server  `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`
}
