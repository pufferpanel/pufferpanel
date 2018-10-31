package oauth2

import (
	"gopkg.in/oauth2.v3"
)

type ClientStore struct {
}

func (cs *ClientStore) GetByID(id string) (oauth2.ClientInfo, error) {
	return &ClientInfo{}, nil
}
