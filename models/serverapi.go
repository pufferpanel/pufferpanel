package models

import (
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/response"
)

type ServerCreation struct {
	pufferpanel.Server

	NodeId uint     `json:"node"`
	Users  []string `json:"users"`
	Name   string   `json:"name"`
}

type GetServerResponse struct {
	Server *ServerView     `json:"server"`
	Perms  *PermissionView `json:"permissions"`
}

type CreateServerResponse struct {
	Id string `json:"id"`
}

type ServerSearchResponse struct {
	Servers []*ServerView `json:"servers"`
	*response.Metadata
}
