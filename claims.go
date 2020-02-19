package pufferpanel

import (
	"crypto/ecdsa"
	"github.com/dgrijalva/jwt-go"
)

type Claim struct {
	jwt.StandardClaims
	PanelClaims PanelClaims `json:"pufferpanel,omitempty"`
}

type PanelClaims struct {
	Scopes map[string][]Scope `json:"scopes"`
}

type Token struct {
	*jwt.Token
	Claims *Claim
}

func ParseToken(publicKey *ecdsa.PublicKey, token string) (*Token, error) {
	claim, err := jwt.ParseWithClaims(token, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	return &Token{
		Token:  claim,
		Claims: claim.Claims.(*Claim),
	}, nil
}
