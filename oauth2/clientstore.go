package oauth2

import (
	"github.com/pufferpanel/pufferpanel/database"
	"gopkg.in/oauth2.v3"
)

type ClientStore struct {
}

func (cs *ClientStore) GetByID(id string) (oauth2.ClientInfo, error) {
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	ci := &ClientInfo{
		ClientID: id,
	}
	res := db.Where(ci).FirstOrInit(ci)

	return ci, res.Error
}

func (cs *ClientStore) Create(id string) (oauth2.ClientInfo, error) {
	return &ClientInfo{}, nil
}