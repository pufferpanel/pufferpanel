package services

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"gorm.io/gorm"
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

// Link replaces the user's link to this issuer
func (oi *OidcIdentity) Link(userId uint, issuer, subject string) error {
	return oi.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.OidcIdentity{}, "user_id = ? AND issuer = ?", userId, issuer).Error; err != nil {
			return err
		}
		return tx.Create(&models.OidcIdentity{UserId: userId, Issuer: issuer, Subject: subject}).Error
	})
}

func (oi *OidcIdentity) GetForUser(userId uint) ([]*models.OidcIdentity, error) {
	identities := make([]*models.OidcIdentity, 0)
	err := oi.DB.Where(&models.OidcIdentity{UserId: userId}).Find(&identities).Error
	return identities, err
}

func (oi *OidcIdentity) DeleteForUser(userId uint) error {
	return oi.DB.Delete(&models.OidcIdentity{}, "user_id = ?", userId).Error
}
