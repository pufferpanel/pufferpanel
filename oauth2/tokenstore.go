package oauth2

import (
	"github.com/pufferpanel/pufferpanel/database"
	"gopkg.in/oauth2.v3"
)

type TokenStore struct {
	//store.TokenStore
}

func (ts *TokenStore) Create(info oauth2.TokenInfo) error {
	db, err := database.GetConnection()

	if err != nil {
		return err
	}

	model := Copy(info)
	return db.Create(model).Error
}

// delete the authorization code
func (ts *TokenStore) RemoveByCode(code string) error {
	return nil
}

// use the access token to delete the token information
func (ts *TokenStore) RemoveByAccess(access string) error {
	obj := &TokenInfo{
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
	obj := &TokenInfo{
		Access: access,
	}

	db, err := database.GetConnection()

	if err != nil {
		return nil, err
	}

	res := db.Where(&obj).FirstOrInit(&obj)
	return obj, res.Error
}

// use the refresh token for token information data
func (ts *TokenStore) GetByRefresh(refresh string) (oauth2.TokenInfo, error) {
	return nil, nil
}
