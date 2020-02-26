package models

import "github.com/pufferpanel/pufferpanel/v2/response"

type UserSearch struct {
	Username  string `form:"username"`
	Email     string `form:"email"`
	PageLimit uint   `form:"limit"`
	Page      uint   `form:"page"`
}

type UserSearchResponse struct {
	Users []*UserView `json:"users"`
	*response.Metadata
}
