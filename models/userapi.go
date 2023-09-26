package models

import (
	"github.com/pufferpanel/pufferpanel/v3"
)

type UserSearch struct {
	Username  string `form:"username"`
	Email     string `form:"email"`
	PageLimit uint   `form:"limit"`
	Page      uint   `form:"page"`
} //@name UserSearch

type UserSearchResponse struct {
	Users []*UserView `json:"users"`
	*pufferpanel.Metadata
} //@name UserSearchResponse
