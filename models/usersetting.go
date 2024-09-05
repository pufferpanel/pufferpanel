package models

type UserSetting struct {
	Key    string `gorm:"column:key;size:100;primaryKey"`
	UserID uint   `gorm:"column:user_id;primaryKey"`
	Value  string `gorm:"column:value;not null;size:4000"`
}
