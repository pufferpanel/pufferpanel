/*
 Copyright 2020 Padduck, LLC
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at
  	http://www.apache.org/licenses/LICENSE-2.0
  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package pufferpanel

import (
	"crypto/ecdsa"
	"errors"
	"github.com/golang-jwt/jwt/v4"
)

type Claim struct {
	jwt.RegisteredClaims
	PanelClaims PanelClaims `json:"pufferpanel,omitempty"`
}

type PanelClaims struct {
	Scopes map[string][]Scope `json:"scopes,omitempty"`
}

type Token struct {
	*jwt.Token
	Claims *Claim
}

func ParseToken(publicKey *ecdsa.PublicKey, token string) (*Token, error) {
	if publicKey == nil {
		return nil, errors.New("PUBLIC KEY NOT LOADED")
	}
	claim, err := jwt.ParseWithClaims(token, &Claim{PanelClaims: PanelClaims{Scopes: make(map[string][]Scope)}}, func(token *jwt.Token) (interface{}, error) {
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
