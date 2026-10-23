package services

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
)

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
	if loaded := oidcProvider.Load(); loaded != nil && loaded.issuer == issuer {
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

// GetOidcProvider returns nil when the issuer cannot be reached; if it could not be reached before, it tries once more
func GetOidcProvider() *oidc.Provider {
	if oidcProvider.Load() == nil {
		LoadOidcProvider()
	}
	if loaded := oidcProvider.Load(); loaded != nil {
		return loaded.provider
	}
	return nil
}
