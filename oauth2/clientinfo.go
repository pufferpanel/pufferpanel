package oauth2

import (
	"github.com/pufferpanel/pufferpanel/models"
	"strconv"
)

type ClientInfo struct {
	ID       uint
	ClientID string
	Secret   string
	UserID   uint
	User     models.User
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
