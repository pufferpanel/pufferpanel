package models

type OAuthClient struct {
	ID int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	UserId int `json:"user_id" db:"user_id"`
	ServerId int `json:"server_id" db:"server_id"`
	Secret string `json:"secret" db:"secret"`
	Description string `json:"description" db:"description"`
}