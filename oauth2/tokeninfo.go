package oauth2

import (
	"github.com/pufferpanel/pufferpanel/models"
	"gopkg.in/oauth2.v3"
	"time"
)

type TokenInfo struct {
	ID uint

	ClientID string
	Client   ClientInfo

	UserID string
	User   models.User

	Scope            string
	Code             string
	CodeCreateAt     time.Time
	CodeExpiresIn    time.Duration
	Access           string
	AccessCreateAt   time.Time
	AccessExpiresIn  time.Duration
	Refresh          string
	RefreshCreateAt  time.Time
	RefreshExpiresIn time.Duration
}

func (ti *TokenInfo) New() oauth2.TokenInfo {
	return &TokenInfo{}
}

func (ti *TokenInfo) GetClientID() string {
	panic("implement me")
}

func (ti *TokenInfo) SetClientID(string) {
	panic("implement me")
}

func (ti *TokenInfo) GetUserID() string {
	panic("implement me")
}

func (ti *TokenInfo) SetUserID(string) {
	panic("implement me")
}

func (ti *TokenInfo) GetRedirectURI() string {
	return ""
}

func (ti *TokenInfo) SetRedirectURI(string) {
}

func (ti *TokenInfo) GetScope() string {
	panic("implement me")
}

func (ti *TokenInfo) SetScope(string) {
	panic("implement me")
}

func (ti *TokenInfo) GetCode() string {
	panic("implement me")
}

func (ti *TokenInfo) SetCode(string) {
	panic("implement me")
}

func (ti *TokenInfo) GetCodeCreateAt() time.Time {
	panic("implement me")
}

func (ti *TokenInfo) SetCodeCreateAt(time.Time) {
	panic("implement me")
}

func (ti *TokenInfo) GetCodeExpiresIn() time.Duration {
	panic("implement me")
}

func (ti *TokenInfo) SetCodeExpiresIn(time.Duration) {
	panic("implement me")
}

func (ti *TokenInfo) GetAccess() string {
	panic("implement me")
}

func (ti *TokenInfo) SetAccess(string) {
	panic("implement me")
}

func (ti *TokenInfo) GetAccessCreateAt() time.Time {
	panic("implement me")
}

func (ti *TokenInfo) SetAccessCreateAt(time.Time) {
	panic("implement me")
}

func (ti *TokenInfo) GetAccessExpiresIn() time.Duration {
	panic("implement me")
}

func (ti *TokenInfo) SetAccessExpiresIn(time.Duration) {
	panic("implement me")
}

func (ti *TokenInfo) GetRefresh() string {
	panic("implement me")
}

func (ti *TokenInfo) SetRefresh(string) {
	panic("implement me")
}

func (ti *TokenInfo) GetRefreshCreateAt() time.Time {
	panic("implement me")
}

func (ti *TokenInfo) SetRefreshCreateAt(time.Time) {
	panic("implement me")
}

func (ti *TokenInfo) GetRefreshExpiresIn() time.Duration {
	panic("implement me")
}

func (ti *TokenInfo) SetRefreshExpiresIn(time.Duration) {
	panic("implement me")
}
