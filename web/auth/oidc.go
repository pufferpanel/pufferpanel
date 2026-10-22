package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"github.com/spf13/cast"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

// The OIDC login only signs in users that already exist in the panel; it never creates them.
// State, nonce and PKCE verifier live in their own short-lived cookie instead of the session: the session cookie
// defaults to SameSite=Strict and is not sent on the redirect back from the identity provider.
const oidcStateCookie = "puffer_oidc"
const oidcStateCookiePath = "/auth/oidc"
const oidcStateMaxAge = 10 * time.Minute

var errOidcNoUser = errors.New("no user matches the identity")
var errOidcEmailNotVerified = errors.New("email of the identity is not verified")

type oidcState struct {
	State    string `json:"s"`
	Nonce    string `json:"n"`
	Verifier string `json:"v"`
}

type oidcUserLookup interface {
	Get(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
}

var oidcProvider *oidc.Provider
var oidcProviderIssuer string
var oidcProviderLock sync.Mutex

// getOidcProvider reads the issuer's discovery document once per issuer (it can be changed in the settings);
// after a failure the next login tries again
func getOidcProvider(ctx context.Context) (*oidc.Provider, error) {
	oidcProviderLock.Lock()
	defer oidcProviderLock.Unlock()

	issuer := config.OidcIssuer.Value()
	if oidcProvider != nil && oidcProviderIssuer == issuer {
		return oidcProvider, nil
	}

	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, err
	}
	oidcProvider = provider
	oidcProviderIssuer = issuer
	return provider, nil
}

func oidcOAuth2Config(provider *oidc.Provider) *oauth2.Config {
	scopeList := strings.Fields(config.OidcScopes.Value())
	if len(scopeList) == 0 {
		scopeList = []string{oidc.ScopeOpenID, "profile", "email"}
	}
	hasOpenId := false
	for _, v := range scopeList {
		if v == oidc.ScopeOpenID {
			hasOpenId = true
		}
	}
	if !hasOpenId {
		scopeList = append([]string{oidc.ScopeOpenID}, scopeList...)
	}

	return &oauth2.Config{
		ClientID:     config.OidcClientId.Value(),
		ClientSecret: config.OidcClientSecret.Value(),
		Endpoint:     provider.Endpoint(),
		RedirectURL:  oidcRedirectUrl(),
		Scopes:       scopeList,
	}
}

func oidcRedirectUrl() string {
	return strings.TrimSuffix(config.MasterUrl.Value(), "/") + "/auth/oidc/callback"
}

func OidcLogin(c *gin.Context) {
	if !config.OidcEnabled.Value() {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	provider, err := getOidcProvider(c.Request.Context())
	if err != nil {
		logging.Error.Printf("OIDC: cannot read the issuer's configuration: %s\n", err)
		oidcFail(c, "unavailable")
		return
	}

	state := oidcState{
		State:    randomOidcValue(),
		Nonce:    randomOidcValue(),
		Verifier: oauth2.GenerateVerifier(),
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
	if err != nil || subtle.ConstantTimeCompare([]byte(state.State), []byte(c.Query("state"))) != 1 {
		oidcFail(c, "state")
		return
	}

	if providerError := c.Query("error"); providerError != "" {
		logging.Info.Printf("OIDC: the identity provider returned %s\n", providerError)
		oidcFail(c, "denied")
		return
	}

	ctx := c.Request.Context()
	provider, err := getOidcProvider(ctx)
	if err != nil {
		logging.Error.Printf("OIDC: cannot read the issuer's configuration: %s\n", err)
		oidcFail(c, "unavailable")
		return
	}

	oauthConfig := oidcOAuth2Config(provider)
	token, err := oauthConfig.Exchange(ctx, c.Query("code"), oauth2.VerifierOption(state.Verifier))
	if err != nil {
		logging.Error.Printf("OIDC: code exchange failed: %s\n", err)
		oidcFail(c, "error")
		return
	}

	rawIdToken, ok := token.Extra("id_token").(string)
	if !ok {
		logging.Error.Printf("OIDC: token response has no id_token\n")
		oidcFail(c, "error")
		return
	}
	idToken, err := provider.Verifier(&oidc.Config{ClientID: oauthConfig.ClientID}).Verify(ctx, rawIdToken)
	if err != nil {
		logging.Error.Printf("OIDC: id_token rejected: %s\n", err)
		oidcFail(c, "error")
		return
	}
	if subtle.ConstantTimeCompare([]byte(idToken.Nonce), []byte(state.Nonce)) != 1 {
		oidcFail(c, "state")
		return
	}

	claims := map[string]interface{}{}
	if err = idToken.Claims(&claims); err != nil {
		oidcFail(c, "error")
		return
	}
	// providers that keep profile and email claims out of the id_token deliver them from the userinfo endpoint
	if _, has := claims[oidcMatchClaim()]; !has {
		userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(token))
		if err == nil && userInfo.Subject == idToken.Subject {
			_ = userInfo.Claims(&claims)
		}
	}

	db := middleware.GetDatabase(c)
	user, err := matchOidcUser(&services.User{DB: db}, claims, config.OidcUserMatch.Value(), oidcUsernameClaim(), config.OidcRequireVerifiedEmail.Value())
	if errors.Is(err, errOidcNoUser) {
		oidcFail(c, "nouser")
		return
	}
	if errors.Is(err, errOidcEmailNotVerified) {
		oidcFail(c, "unverified")
		return
	}
	if err != nil {
		logging.Error.Printf("OIDC: user lookup failed: %s\n", err)
		oidcFail(c, "error")
		return
	}

	if _, ok = startSession(c, user); !ok {
		return
	}
	// the frontend loads the scopes there before it renders, like after a password login
	c.Redirect(http.StatusFound, "/auth/oidc/complete")
}

func oidcMatchClaim() string {
	if strings.EqualFold(config.OidcUserMatch.Value(), "username") {
		return oidcUsernameClaim()
	}
	return "email"
}

// oidcUsernameClaim falls back to preferred_username when the setting was saved empty
func oidcUsernameClaim() string {
	if claim := config.OidcUsernameClaim.Value(); claim != "" {
		return claim
	}
	return "preferred_username"
}

// matchOidcUser finds the panel user for the identity: by email (default) or by the configured username claim
func matchOidcUser(lookup oidcUserLookup, claims map[string]interface{}, mode, usernameClaim string, requireVerifiedEmail bool) (*models.User, error) {
	var user *models.User
	var err error

	if strings.EqualFold(mode, "username") {
		username, _ := claims[usernameClaim].(string)
		if username == "" {
			return nil, errOidcNoUser
		}
		user, err = lookup.Get(username)
	} else {
		email, _ := claims["email"].(string)
		if email == "" {
			return nil, errOidcNoUser
		}
		if requireVerifiedEmail && !cast.ToBool(claims["email_verified"]) {
			return nil, errOidcEmailNotVerified
		}
		user, err = lookup.GetByEmail(email)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errOidcNoUser
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func oidcFail(c *gin.Context, reason string) {
	c.Redirect(http.StatusFound, "/auth/login?oidcError="+url.QueryEscape(reason))
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
