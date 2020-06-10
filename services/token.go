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

package services

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var signingMethod = jwt.SigningMethodES256
var privateKey *ecdsa.PrivateKey
var publicKey *ecdsa.PublicKey
var locker sync.Mutex
var timer time.Time

func GetPublicKey() *ecdsa.PublicKey {
	ValidateTokenLoaded()
	if privateKey != nil {
		return &privateKey.PublicKey
	} else {
		return publicKey
	}
}

func Generate(claims jwt.Claims) (string, error) {
	ValidateTokenLoaded()
	token := jwt.NewWithClaims(signingMethod, claims)
	return token.SignedString(privateKey)
}

func GenerateSession(id uint) (string, error) {
	claims := &pufferpanel.Claim{
		StandardClaims: jwt.StandardClaims{
			Audience:  "session",
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   strconv.Itoa(int(id)),
		},
	}

	return Generate(claims)
}

func GenerateOAuthForClient(client *models.Client) (string, error) {
	claims := &pufferpanel.Claim{
		StandardClaims: jwt.StandardClaims{
			Audience:  "oauth2",
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		PanelClaims: pufferpanel.PanelClaims{
			Scopes: map[string][]pufferpanel.Scope{
				client.ServerId: client.Scopes,
			},
		},
	}

	return Generate(claims)
}

func GenerateOAuthForNode(nodeId uint) (string, error) {
	claims := &pufferpanel.Claim{
		StandardClaims: jwt.StandardClaims{
			Audience:  "oauth2",
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		PanelClaims: pufferpanel.PanelClaims{
			Scopes: map[string][]pufferpanel.Scope{
				"": {pufferpanel.ScopeOAuth2Auth},
			},
		},
	}
	return Generate(claims)
}

func (ps *Permission) GenerateOAuthForUser(userId uint, serverId *string) (string, error) {
	var err error
	var permissions []*models.Permissions

	if serverId == nil {
		permissions, err = ps.GetForUser(userId)
	} else {
		var perm *models.Permissions
		perm, err = ps.GetForUserAndServer(userId, serverId)
		if err != nil {
			return "", err
		}
		permissions = []*models.Permissions{perm}

		perm, err = ps.GetForUserAndServer(userId, nil)
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			return "", err
		}
		if perm.ID != 0 && !gorm.IsRecordNotFoundError(err) {
			permissions = append(permissions, perm)
		}
	}

	if err != nil {
		return "", err
	}

	claims := &pufferpanel.Claim{
		StandardClaims: jwt.StandardClaims{
			Audience:  "oauth2",
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		PanelClaims: pufferpanel.PanelClaims{
			Scopes: map[string][]pufferpanel.Scope{},
		},
	}

	for _, perm := range permissions {
		var existing []pufferpanel.Scope
		if perm.ServerIdentifier == nil {
			existing = claims.PanelClaims.Scopes[""]
		} else {
			existing = claims.PanelClaims.Scopes[*perm.ServerIdentifier]
		}

		if existing == nil {
			existing = make([]pufferpanel.Scope, 0)
		}

		existing = append(existing, perm.ToScopes()...)

		if perm.ServerIdentifier == nil {
			claims.PanelClaims.Scopes[""] = existing
		} else {
			claims.PanelClaims.Scopes[*perm.ServerIdentifier] = existing
		}
	}

	return Generate(claims)
}

func ParseToken(token string) (*pufferpanel.Token, error) {
	ValidateTokenLoaded()
	return pufferpanel.ParseToken(publicKey, token)
}

func ValidateTokenLoaded() {
	locker.Lock()
	defer locker.Unlock()
	//only load public if panel is disabled
	if !viper.GetBool("panel.enable") {
		if publicKey == nil || timer.Before(time.Now()) {
			loadPublic()
		}
	} else if privateKey == nil {
		loadPrivate()
	}
}

func loadPrivate() {
	var privKey *ecdsa.PrivateKey
	privKeyFile, err := os.OpenFile(viper.GetString("token.private"), os.O_RDONLY, 0600)
	defer pufferpanel.Close(privKeyFile)
	if os.IsNotExist(err) {
		privKey, err = generatePrivateKey()
	} else if err == nil {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, privKeyFile)
		block, _ := pem.Decode(buf.Bytes())

		privKey, err = ecdsa.GenerateKey(elliptic.P256(), bytes.NewReader(block.Bytes))
	}

	if err != nil {
		logging.Error().Printf("Internal error on token service: %s", err)
		return
	}

	privateKey = privKey
	publicKey = &privateKey.PublicKey
	return
}

func generatePrivateKey() (privKey *ecdsa.PrivateKey, err error) {
	var key bytes.Buffer
	privKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return
	}

	privKeyEncoded, _ := x509.MarshalECPrivateKey(privKey)
	privKeyFile, err := os.OpenFile(viper.GetString("token.private"), os.O_CREATE|os.O_WRONLY, 0600)
	defer pufferpanel.Close(privKeyFile)
	if err != nil {
		return
	}
	err = pem.Encode(privKeyFile, &pem.Block{Type: "PRIVATE KEY", Bytes: privKeyEncoded})
	if err != nil {
		return
	}
	err = pem.Encode(&key, &pem.Block{Type: "PRIVATE KEY", Bytes: privKeyEncoded})
	if err != nil {
		return
	}

	return
}

func loadPublic() {
	pubKeyPath := viper.GetString("token.public")
	var pubKeyFile io.ReadCloser
	var err error

	if strings.HasPrefix("https://", pubKeyPath) || strings.HasPrefix("http://", pubKeyPath) {
		client := http.Client{}
		response, err := client.Get(pubKeyPath)
		if err != nil {
			logging.Error().Printf("Internal error on token service: %s", err)
			return
		}
		pubKeyFile = response.Body
		timer = time.Now().Add(5 * time.Minute)
	} else {
		pubKeyFile, err = os.OpenFile(pubKeyPath, os.O_RDONLY, 0644)
		if err != nil {
			logging.Error().Printf("Internal error on token service: %s", err)
			return
		}
	}

	defer pufferpanel.Close(pubKeyFile)

	data, err := ioutil.ReadAll(pubKeyFile)
	if err != nil {
		logging.Error().Printf("Internal error on token service: %s", err)
		return
	}

	block, _ := pem.Decode(data)
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		logging.Error().Printf("Internal error on token service: %s", err)
		return
	}

	var ok bool
	publicKey, ok = key.(*ecdsa.PublicKey)
	if !ok {
		logging.Error().Printf("Internal error on token service: %s", errors.New("public key is not ECDSA"))
		return
	}
}