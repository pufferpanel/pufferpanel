package models

import (
	"strconv"
)

type ClientInfo struct {
	ID       uint   `json:"-"`
	ClientID string `gorm:"NOT NULL" json:"-"`
	Secret   string `gorm:"NOT NULL" json:"-"`
	UserID   uint   `json:"-"`
	User     User   `gorm:"save_associations:false" json:"-"`
	Panel    bool   `gorm:"NOT NULL; DEFAULT:0"`

	ServerScopes []ClientServerScopes `gorm:"save_associations:false preload:false" json:"-"`
}

type ClientInfos []*ClientInfo

type ClientServerScopes struct {
	ID           uint   `json:"-"`
	ClientInfoID uint   `gorm:"NOT NULL" json:"-"`
	ServerId     *uint  `json:"-"`
	Server       Server `gorm:"save_associations:false" json:"-"`
	Scope        string `gorm:"NOT NULL" json:"-"`
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

func (ci *ClientInfo) MergeServers() (map[string][]string, []string) {
	mapping := make(map[string][]string, 0)
	plain := make([]string, 0)

	for _, v := range ci.ServerScopes {
		temp := mapping[v.Server.Identifier]
		if temp == nil {
			temp = make([]string, 0)
		}
		temp = append(temp, v.Scope)

		mapping[v.Server.Identifier] = temp
		plain = append(plain, v.Server.Identifier+":"+v.Scope)
	}

	return mapping, plain
}
