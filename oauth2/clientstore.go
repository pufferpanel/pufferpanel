package oauth2

import (
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"gopkg.in/oauth2.v3"
)

type ClientStore struct {
}

func (cs *ClientStore) GetByID(id string) (oauth2.ClientInfo, error) {
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	ci := &models.ClientInfo{
		ClientID: id,
	}
	res := db.Where(ci).FirstOrInit(ci)

	return ci, res.Error
}

func (cs *ClientStore) Create(id string) (oauth2.ClientInfo, error) {
	return &models.ClientInfo{}, nil
}