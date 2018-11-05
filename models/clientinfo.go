package models

import (
	"strconv"
)

type ClientInfo struct {
	ID       uint
	ClientID string
	Secret   string
	UserID   uint
	User     User

	ServerScopes []ClientServerScopes
}

type ClientServerScopes struct {
	ID           uint
	ClientInfoID uint
	ServerId     uint
	Server       Server
	Scope        string
}

func (ci *ClientInfo) GetSecret() string {
	return ci.Secret
}

func (ci *ClientInfo) GetID() string {
	return ci.ClientID
}

func (ci *ClientInfo) GetDomain() string {
	return "*"
}

func (ci *ClientInfo) GetUserID() string {
	return strconv.Itoa(int(ci.UserID))
}

func (ci *ClientInfo) MergeServers() map[string][]string {
	mapping := make(map[string][]string, 0)

	for _, v := range ci.ServerScopes {
		temp := mapping[v.Server.Identifier]
		if temp == nil {
			temp = make([]string, 0)
		}
		temp = append(temp, v.Scope)

		mapping[v.Server.Identifier] = temp
	}

	return mapping
}
