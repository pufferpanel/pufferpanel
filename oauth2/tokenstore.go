package oauth2

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/models"
	"gopkg.in/oauth2.v3"
	"time"
)

type TokenStore struct {
	//store.TokenStore
}

func (ts *TokenStore) Create(info oauth2.TokenInfo) error {
	db, err := database.GetConnection()

	if err != nil {
		return err
	}

	model := models.Copy(info)

	if model.ClientInfoID == 0 {
		client := &models.ClientInfo{
			ClientID: info.GetClientID(),
		}

		err = db.Where(client).First(client).Error
		if err != nil {
			return err
		}
		model.ClientInfoID = client.ID
		model.ClientInfo = *client
	}

	return db.Create(model).Error
}

// delete the authorization code
func (ts *TokenStore) RemoveByCode(code string) error {
	return nil
}

// use the access token to delete the token information
func (ts *TokenStore) RemoveByAccess(access string) error {
	obj := &models.TokenInfo{
		Access: access,
	}

	db, err := database.GetConnection()

	if err != nil {
		return err
	}

	return db.Delete(obj).Error
}

// use the refresh token to delete the token information
func (ts *TokenStore) RemoveByRefresh(refresh string) error {
	return nil

}

// use the authorization code for token information data
func (ts *TokenStore) GetByCode(code string) (oauth2.TokenInfo, error) {
	return nil, nil
}

// use the access token for token information data
func (ts *TokenStore) GetByAccess(access string) (oauth2.TokenInfo, error) {
	obj := &models.TokenInfo{
		Access: access,
	}

	db, err := database.GetConnection()

	if err != nil {
		return nil, err
	}

	res := db.Preload("ClientInfo").Where(obj).First(obj)
	err = res.Error

	if obj == nil || obj.ID == 0 || err != nil {
		if err == nil || gorm.IsRecordNotFoundError(err) {
			err = errors.New("token is invalid")
		}
		return nil, err
	}

	if obj.GetAccessCreateAt().Add(obj.GetAccessExpiresIn()).Before(time.Now()) {
		return nil, err
	}

	db.Preload("ServerScopes").Preload("User").Where(&obj.ClientInfo).First(&obj.ClientInfo)
	return obj, err
}

// use the refresh token for token information data
func (ts *TokenStore) GetByRefresh(refresh string) (oauth2.TokenInfo, error) {
	return nil, nil
}
