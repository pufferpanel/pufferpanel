package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/scopes"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

// The OIDC login only signs in users that already exist in the panel and linked their provider account; it never
// creates users. State, nonce and PKCE verifier live in their own short-lived cookie: gin's session cookie is
// SameSite=Strict and is not sent on the redirect back from the identity provider.
const oidcStateCookie = "puffer_oidc"
const oidcStateCookiePath = "/auth/oidc"
const oidcStateMaxAge = 10 * time.Minute

type oidcState struct {
	State    string `json:"s"`
	Nonce    string `json:"n"`
	Verifier string `json:"v"`
	Link     bool   `json:"l,omitempty"`
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

func OidcLogin(c *gin.Context) {
	if !config.OidcEnabled.Value() {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	provider := services.GetOidcProvider()
	if provider == nil {
		oidcFail(c, "unavailable")
		return
	}

	state := oidcState{
		State:    randomOidcValue(),
		Nonce:    randomOidcValue(),
		Verifier: oauth2.GenerateVerifier(),
		Link:     c.Query("link") == "1",
	}
	data, err := json.Marshal(state)
	if err != nil {
		oidcFail(c, "error")
		return
	}
	setOidcStateCookie(c, base64.RawURLEncoding.EncodeToString(data), int(oidcStateMaxAge/time.Second))

	authUrl := oidcOAuth2Config(provider).AuthCodeURL(state.State, oidc.Nonce(state.Nonce), oauth2.S256ChallengeOption(state.Verifier))
	c.Redirect(http.StatusFound, authUrl)
}

func OidcCallback(c *gin.Context) {
	if !config.OidcEnabled.Value() {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	rawState, err := c.Cookie(oidcStateCookie)
	// one attempt per state, whatever happens next
	setOidcStateCookie(c, "", -1)
	if err != nil {
		oidcFail(c, "state")
		return
	}
	state, err := decodeOidcState(rawState)
	if err != nil || !oidcValuesMatch(state.State, c.Query("state")) {
		oidcFail(c, "state")
		return
	}
	fail := oidcFailFor(c, state)

	if providerError := c.Query("error"); providerError != "" {
		logging.Info.Printf("OIDC: the identity provider returned %s\n", providerError)
		fail("denied")
		return
	}

	ctx := c.Request.Context()
	provider := services.GetOidcProvider()
	if provider == nil {
		fail("unavailable")
		return
	}

	oauthConfig := oidcOAuth2Config(provider)
	token, err := oauthConfig.Exchange(ctx, c.Query("code"), oauth2.VerifierOption(state.Verifier))
	if err != nil {
		logging.Error.Printf("OIDC: code exchange failed: %s\n", err)
		fail("error")
		return
	}

	rawIdToken, ok := token.Extra("id_token").(string)
	if !ok {
		logging.Error.Printf("OIDC: token response has no id_token\n")
		fail("error")
		return
	}
	idToken, err := provider.Verifier(&oidc.Config{ClientID: oauthConfig.ClientID}).Verify(ctx, rawIdToken)
	if err != nil {
		logging.Error.Printf("OIDC: id_token rejected: %s\n", err)
		fail("error")
		return
	}
	if !oidcValuesMatch(idToken.Nonce, state.Nonce) {
		fail("state")
		return
	}

	db := middleware.GetDatabase(c)
	oi := &services.OidcIdentity{DB: db}
	if state.Link {
		linkOidcIdentity(c, oi, idToken, fail)
		return
	}

	user, err := oi.GetUser(idToken.Issuer, idToken.Subject)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		oidcFail(c, "notlinked")
		return
	}
	if err != nil {
		logging.Error.Printf("OIDC: user lookup failed: %s\n", err)
		oidcFail(c, "error")
		return
	}

	userScopes, _, err := startSession(c, user)
	if errors.Is(err, pufferpanel.ErrLoginNotPermitted) {
		oidcFail(c, "nologin")
		return
	} else if err != nil {
		logging.Error.Printf("OIDC: cannot start the session: %s\n", err)
		oidcFail(c, "error")
		return
	}

	if err = setOidcScopesCookie(c, userScopes); err != nil {
		oidcFail(c, "error")
		return
	}
	c.Redirect(http.StatusFound, "/")
}

// linkOidcIdentity links the provider account to the user who is signed in to the panel in this browser; the login
// cookie puffer_auth has no SameSite attribute, so it arrives with the provider's redirect
func linkOidcIdentity(c *gin.Context, oi *services.OidcIdentity, idToken *oidc.IDToken, fail func(string)) {
	user := oidcSessionUser(c)
	if user == nil {
		fail("error")
		return
	}

	linked, err := oi.GetUser(idToken.Issuer, idToken.Subject)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		// not linked to anyone yet
	case err != nil:
		logging.Error.Printf("OIDC: user lookup failed: %s\n", err)
		fail("error")
		return
	case linked.ID != user.ID:
		fail("taken")
		return
	}

	if err = oi.Link(user.ID, idToken.Issuer, idToken.Subject); err != nil {
		logging.Error.Printf("OIDC: cannot link the identity to user %d: %s\n", user.ID, err)
		fail("error")
		return
	}
	c.Redirect(http.StatusFound, "/self?oidc=linked#security")
}

// oidcSessionUser returns the user of this browser's panel session if they may edit their own account
func oidcSessionUser(c *gin.Context) *models.User {
	token, err := c.Cookie("puffer_auth")
	if err != nil || token == "" {
		return nil
	}
	db := middleware.GetDatabase(c)
	session, err := (&services.Session{DB: db}).Validate(token)
	if err != nil || !isUserSession(session) {
		return nil
	}
	perms, err := (&services.Permission{DB: db}).GetForUserAndServer(session.User.ID, "")
	if err != nil || !scopes.ContainsScope(perms.Scopes, scopes.ScopeSelfEdit) {
		return nil
	}
	return &session.User
}

// isUserSession is false for sessions of OAuth2 clients, which belong to no user
func isUserSession(session *models.Session) bool {
	return session.UserId != nil
}

// oidcValuesMatch compares state and nonce values in constant time
func oidcValuesMatch(expected, actual string) bool {
	return subtle.ConstantTimeCompare([]byte(expected), []byte(actual)) == 1
}

func oidcFail(c *gin.Context, reason string) {
	c.Redirect(http.StatusFound, "/auth/oidc/error/"+url.PathEscape(reason))
}

// oidcFailFor sends a failed login to the error page and a failed link back to the security settings, where the
// user is still signed in
func oidcFailFor(c *gin.Context, state oidcState) func(reason string) {
	if !state.Link {
		return func(reason string) { oidcFail(c, reason) }
	}
	return func(reason string) {
		c.Redirect(http.StatusFound, "/self?oidc="+url.QueryEscape(reason)+"#security")
	}
}

// setOidcScopesCookie stores the scopes where the frontend keeps them after a password login, so the first page
// after the redirect already knows the user's permissions
func setOidcScopesCookie(c *gin.Context, userScopes []*scopes.Scope) error {
	data, err := json.Marshal(userScopes)
	if err != nil {
		return err
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "puffer_scopes",
		Value:    url.QueryEscape(string(data)),
		Path:     "/",
		Secure:   c.Request.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}

func setOidcStateCookie(c *gin.Context, value string, maxAge int) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     oidcStateCookie,
		Value:    value,
		Path:     oidcStateCookiePath,
		MaxAge:   maxAge,
		Secure:   c.Request.TLS != nil,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func decodeOidcState(raw string) (oidcState, error) {
	var state oidcState
	data, err := base64.RawURLEncoding.DecodeString(raw)
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

func randomOidcValue() string {
	data := make([]byte, 32)
	_, _ = rand.Read(data)
	return base64.RawURLEncoding.EncodeToString(data)
}
