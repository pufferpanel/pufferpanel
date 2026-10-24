package services

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/scopes"
	"github.com/pufferpanel/pufferpanel/v3/utils"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

var ErrOidcUnavailable = errors.New("the OIDC provider cannot be reached")
var ErrOidcStateMismatch = errors.New("the OIDC state or nonce does not match")
var ErrOidcIdentityTaken = errors.New("the OIDC identity is linked to another user")

type loadedOidcProvider struct {
	issuer   string
	provider *oidc.Provider
}

var oidcProvider atomic.Pointer[loadedOidcProvider]

// LoadOidcProvider reads the issuer's discovery document; it runs at startup and after the settings are saved,
// and only fetches again when the issuer changed or the last attempt failed
func LoadOidcProvider() {
	issuer := config.OidcIssuer.Value()
	if !config.OidcEnabled.Value() || issuer == "" {
		oidcProvider.Store(nil)
		return
	}
	if oidcProviderLoadedFor(issuer) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		logging.Error.Printf("OIDC: cannot read the issuer's configuration: %s\n", err)
		oidcProvider.Store(nil)
		return
	}
	oidcProvider.Store(&loadedOidcProvider{issuer: issuer, provider: provider})
}

func oidcProviderLoadedFor(issuer string) bool {
	loaded := oidcProvider.Load()
	return loaded != nil && loaded.issuer == issuer
}

// availableOidcProvider returns ErrOidcUnavailable when the issuer cannot be reached; if it could not be reached
// before, it tries once more
func availableOidcProvider() (*oidc.Provider, error) {
	if oidcProvider.Load() == nil {
		LoadOidcProvider()
	}
	if loaded := oidcProvider.Load(); loaded != nil {
		return loaded.provider, nil
	}
	return nil, ErrOidcUnavailable
}

func oidcOAuth2Config(provider *oidc.Provider) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.OidcClientId.Value(),
		ClientSecret: config.OidcClientSecret.Value(),
		Endpoint:     provider.Endpoint(),
		RedirectURL:  oidcRedirectUrl(),
		Scopes:       oidcScopes(config.OidcScopes.Value()),
	}
}

// oidcScopes splits the space-separated setting; openid is always requested, the protocol requires it
func oidcScopes(value string) []string {
	scopeList := strings.Fields(value)
	if !slices.Contains(scopeList, oidc.ScopeOpenID) {
		scopeList = append([]string{oidc.ScopeOpenID}, scopeList...)
	}
	return scopeList
}

func oidcRedirectUrl() string {
	return strings.TrimSuffix(config.MasterUrl.Value(), "/") + "/auth/oidc/callback"
}

// OidcState is what a login remembers between the redirect to the provider and the callback
type OidcState struct {
	State    string `json:"s"`
	Nonce    string `json:"n"`
	Verifier string `json:"v"`
	Link     bool   `json:"l,omitempty"`
}

func NewOidcState(link bool) (OidcState, error) {
	state, err := utils.GenerateRandomString(32)
	if err != nil {
		return OidcState{}, err
	}
	nonce, err := utils.GenerateRandomString(32)
	if err != nil {
		return OidcState{}, err
	}
	return OidcState{State: state, Nonce: nonce, Verifier: oauth2.GenerateVerifier(), Link: link}, nil
}

func (s OidcState) Encode() (string, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(data), nil
}

func DecodeOidcState(raw string) (OidcState, error) {
	var state OidcState
	data, err := base64.URLEncoding.DecodeString(raw)
	if err != nil {
		return state, err
	}
	if err = json.Unmarshal(data, &state); err != nil {
		return state, err
	}
	if state.State == "" || state.Nonce == "" || state.Verifier == "" {
		return state, errors.New("incomplete oidc state")
	}
	return state, nil
}

// Matches compares the state the provider sent back
func (s OidcState) Matches(state string) bool {
	return oidcValuesMatch(s.State, state)
}

// oidcValuesMatch compares state and nonce values in constant time
func oidcValuesMatch(expected, actual string) bool {
	return subtle.ConstantTimeCompare([]byte(expected), []byte(actual)) == 1
}

// OidcAuthCodeURL is where the browser signs in at the provider
func OidcAuthCodeURL(state OidcState) (string, error) {
	provider, err := availableOidcProvider()
	if err != nil {
		return "", err
	}
	return oidcOAuth2Config(provider).AuthCodeURL(state.State, oidc.Nonce(state.Nonce), oauth2.S256ChallengeOption(state.Verifier)), nil
}

// OidcIdentityFromCode exchanges the code from the callback and returns the verified provider account
func OidcIdentityFromCode(ctx context.Context, code string, state OidcState) (issuer, subject string, err error) {
	provider, err := availableOidcProvider()
	if err != nil {
		return "", "", err
	}

	oauthConfig := oidcOAuth2Config(provider)
	token, err := oauthConfig.Exchange(ctx, code, oauth2.VerifierOption(state.Verifier))
	if err != nil {
		return "", "", fmt.Errorf("code exchange failed: %w", err)
	}
	rawIdToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", "", errors.New("token response has no id_token")
	}
	idToken, err := provider.Verifier(&oidc.Config{ClientID: oauthConfig.ClientID}).Verify(ctx, rawIdToken)
	if err != nil {
		return "", "", fmt.Errorf("id_token rejected: %w", err)
	}
	if !oidcValuesMatch(idToken.Nonce, state.Nonce) {
		return "", "", ErrOidcStateMismatch
	}
	return idToken.Issuer, idToken.Subject, nil
}

type OidcIdentity struct {
	DB *gorm.DB
}

// GetUser returns the user linked to the provider account, or gorm.ErrRecordNotFound
func (oi *OidcIdentity) GetUser(issuer, subject string) (*models.User, error) {
	identity := &models.OidcIdentity{}
	err := oi.DB.Preload("User").Where(&models.OidcIdentity{Issuer: issuer, Subject: subject}).First(identity).Error
	if err != nil {
		return nil, err
	}
	return &identity.User, nil
}

// Link replaces the user's link to this issuer; a provider account of another user returns ErrOidcIdentityTaken
func (oi *OidcIdentity) Link(userId uint, issuer, subject string) error {
	linked, err := oi.GetUser(issuer, subject)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		// not linked to anyone yet
	case err != nil:
		return err
	case linked.ID != userId:
		return ErrOidcIdentityTaken
	}

	return oi.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.OidcIdentity{}, "user_id = ? AND issuer = ?", userId, issuer).Error; err != nil {
			return err
		}
		return tx.Create(&models.OidcIdentity{UserId: userId, Issuer: issuer, Subject: subject}).Error
	})
}

// LinkingUser returns the user of a panel session if they may edit their own account, otherwise nil
func (oi *OidcIdentity) LinkingUser(sessionToken string) *models.User {
	if sessionToken == "" {
		return nil
	}
	session, err := (&Session{DB: oi.DB}).Validate(sessionToken)
	if err != nil || !session.IsUser() {
		return nil
	}
	perms, err := (&Permission{DB: oi.DB}).GetForUserAndServer(session.User.ID, "")
	if err != nil || !scopes.ContainsScope(perms.Scopes, scopes.ScopeSelfEdit) {
		return nil
	}
	return &session.User
}

func (oi *OidcIdentity) GetForUser(userId uint) ([]*models.OidcIdentity, error) {
	identities := make([]*models.OidcIdentity, 0)
	err := oi.DB.Where(&models.OidcIdentity{UserId: userId}).Find(&identities).Error
	return identities, err
}

func (oi *OidcIdentity) DeleteForUser(userId uint) error {
	return oi.DB.Delete(&models.OidcIdentity{}, "user_id = ?", userId).Error
}
