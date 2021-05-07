package models

type UserSetting struct {
	Key    string `gorm:"NOT NULL;size:100;PRIMARY_KEY"`
	UserID uint   `gorm:"NOT NULL;PRIMARY_KEY;AUTO_INCREMENT:false"`
	Value  string `gorm:"NOT NULL;type:text"`
}

type UserSettings []*UserSetting
