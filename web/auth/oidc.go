package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/scopes"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"gorm.io/gorm"
)

// The OIDC login only signs in users that already exist in the panel and linked their provider account; it never
// creates users. State, nonce and PKCE verifier live in their own short-lived cookie: gin's session cookie is
// SameSite=Strict and is not sent on the redirect back from the identity provider.
const oidcStateCookie = "puffer_oidc"
const oidcStateCookiePath = "/auth/oidc"
const oidcStateMaxAge = 10 * time.Minute

func OidcLogin(c *gin.Context) {
	if !config.OidcEnabled.Value() {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	state, err := services.NewOidcState(c.Query("link") == "1")
	if err != nil {
		oidcFail(c, "error")
		return
	}
	authUrl, err := services.OidcAuthCodeURL(state)
	if err != nil {
		oidcFail(c, "unavailable")
		return
	}
	encodedState, err := state.Encode()
	if err != nil {
		oidcFail(c, "error")
		return
	}
	setOidcStateCookie(c, encodedState, int(oidcStateMaxAge/time.Second))
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
	state, err := services.DecodeOidcState(rawState)
	if err != nil || !state.Matches(c.Query("state")) {
		oidcFail(c, "state")
		return
	}
	fail := oidcFailFor(c, state)

	if providerError := c.Query("error"); providerError != "" {
		logging.Error.Printf("OIDC: the identity provider returned %s\n", providerError)
		fail("denied")
		return
	}

	issuer, subject, err := services.OidcIdentityFromCode(c.Request.Context(), c.Query("code"), state)
	switch {
	case errors.Is(err, services.ErrOidcUnavailable):
		fail("unavailable")
		return
	case errors.Is(err, services.ErrOidcStateMismatch):
		fail("state")
		return
	case err != nil:
		logging.Error.Printf("OIDC: %s\n", err)
		fail("error")
		return
	}

	oi := &services.OidcIdentity{DB: middleware.GetDatabase(c)}
	if state.Link {
		linkOidcIdentity(c, oi, issuer, subject, fail)
		return
	}

	user, err := oi.GetUser(issuer, subject)
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
func linkOidcIdentity(c *gin.Context, oi *services.OidcIdentity, issuer, subject string, fail func(string)) {
	user := oi.LinkingUser(middleware.GetToken(c))
	if user == nil {
		fail("error")
		return
	}

	err := oi.Link(user.ID, issuer, subject)
	if errors.Is(err, services.ErrOidcIdentityTaken) {
		fail("taken")
		return
	}
	if err != nil {
		logging.Error.Printf("OIDC: cannot link the identity to user %d: %s\n", user.ID, err)
		fail("error")
		return
	}
	c.Redirect(http.StatusFound, "/self?oidc=linked#security")
}

func oidcFail(c *gin.Context, reason string) {
	c.Redirect(http.StatusFound, "/auth/oidc/error/"+url.PathEscape(reason))
}

// oidcFailFor sends a failed login to the error page and a failed link back to the security settings, where the
// user is still signed in
func oidcFailFor(c *gin.Context, state services.OidcState) func(reason string) {
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
