package oauth2

import (
	"gopkg.in/oauth2.v3"
)

type TokenStore struct {
	//store.TokenStore
}

func (ts *TokenStore) Create(info oauth2.TokenInfo) error {
	return nil
}

// delete the authorization code
func (ts *TokenStore) RemoveByCode(code string) error {
	return nil
}

// use the access token to delete the token information
func (ts *TokenStore) RemoveByAccess(access string) error {
	return nil
}

// use the refresh token to delete the token information
func (ts *TokenStore) RemoveByRefresh(refresh string) error {
	return nil

}

// use the authorization code for token information data
func (ts *TokenStore) GetByCode(code string) (oauth2.TokenInfo, error) {
	return &TokenInfo{}, nil
}

// use the access token for token information data
func (ts *TokenStore) GetByAccess(access string) (oauth2.TokenInfo, error) {
	return &TokenInfo{}, nil
}

// use the refresh token for token information data
func (ts *TokenStore) GetByRefresh(refresh string) (oauth2.TokenInfo, error) {
	return &TokenInfo{}, nil
}
