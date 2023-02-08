package models

type UserSetting struct {
	Key    string `gorm:"NOT NULL;size:100;primaryKey"`
	UserID uint   `gorm:"NOT NULL;primaryKey;AUTO_INCREMENT:false"`
	Value  string `gorm:"NOT NULL;type:text"`
}
