package models

import (
	"github.com/pufferpanel/pufferpanel/v3"
)

type ServerCreation struct {
	pufferpanel.Server

	NodeId uint     `json:"node"`
	Users  []string `json:"users"`
	Name   string   `json:"name"`
} //@name CreatedServer

type GetServerResponse struct {
	Server *ServerView     `json:"server"`
	Perms  *PermissionView `json:"permissions"`
} //@name GetServer

type CreateServerResponse struct {
	Id string `json:"id"`
} //@name CreatedServerId

type ServerSearchResponse struct {
	Servers []*ServerView `json:"servers"`
	*pufferpanel.Metadata
} //@name ServerSearchResults

type ServerWithName struct {
	pufferpanel.Server
	Name string `json:"name"`
} //@name NamedServer
