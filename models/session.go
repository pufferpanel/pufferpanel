package models

import (
	"time"
)

type Session struct {
	ID             uint      `gorm:"column:id;primaryKey;autoIncrement" json:"-"`
	Token          string    `gorm:"column:token;not null;size:64;uniqueIndex;unique" json:"-"`
	ExpirationTime time.Time `gorm:"column:expiration_time;not null;index" json:"-"`

	UserId *uint `gorm:"column:user_id;index" json:"-"`
	User   User  `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

	ClientId *uint  `gorm:"column:client_id;index" json:"-"`
	Client   Client `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

	//if this set is for a server, what server
	ServerIdentifier *string `gorm:"column:server_identifier" json:"-"`
	Server           Server  `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`
}
