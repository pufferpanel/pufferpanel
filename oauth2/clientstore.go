package oauth2

import (
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/store"
)

type ClientStore struct {
	store.ClientStore
}

func (cs *ClientStore) GetByID(id string) (oauth2.ClientInfo, error) {
	return &ClientInfo{}, nil
}
