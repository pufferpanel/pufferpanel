package models

import (
	"time"
)

// OidcIdentity links a panel user to the account of an OpenID Connect provider (its issuer and the stable sub claim)
type OidcIdentity struct {
	ID uint `gorm:"column:id;primaryKey;autoIncrement" json:"-"`

	UserId uint `gorm:"column:user_id;not null;index" json:"-"`
	User   User `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

	Issuer  string `gorm:"column:issuer;not null;size:255;uniqueIndex:idx_oidc_identity" json:"issuer"`
	Subject string `gorm:"column:subject;not null;size:255;uniqueIndex:idx_oidc_identity" json:"-"`

	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
} //@name OidcIdentity
