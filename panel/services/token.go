package services

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/daemon"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/panel/models"
	"github.com/pufferpanel/pufferpanel/v2/scope"
	"github.com/spf13/viper"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

var signingMethod = jwt.SigningMethodES256
var privateKey *ecdsa.PrivateKey
var locker sync.Mutex

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
			Scopes: map[string][]scope.Scope{
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
			Scopes: map[string][]scope.Scope{
				"": {scope.OAuth2Auth},
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
			Scopes: map[string][]scope.Scope{},
		},
	}

	for _, perm := range permissions {
		var existing []scope.Scope
		if perm.ServerIdentifier == nil {
			existing = claims.PanelClaims.Scopes[""]
		} else {
			existing = claims.PanelClaims.Scopes[*perm.ServerIdentifier]
		}

		if existing == nil {
			existing = make([]scope.Scope, 0)
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
	return pufferpanel.ParseToken(&privateKey.PublicKey, token)
}

func ValidateTokenLoaded() {
	locker.Lock()
	defer locker.Unlock()
	if privateKey == nil {
		load()
	}
}

func load() {
	var privKey *ecdsa.PrivateKey
	privKeyFile, err := os.OpenFile(viper.GetString("panel.token.private"), os.O_RDONLY, 0600)
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

	pubKey := &privateKey.PublicKey
	pubKeyEncoded, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		logging.Error().Printf("Internal error on token service: %s", err)
		return
	}

	pubKeyFile, err := os.OpenFile(viper.GetString("panel.token.public"), os.O_CREATE|os.O_RDWR, 0644)
	defer pufferpanel.Close(pubKeyFile)
	if err != nil {
		logging.Error().Printf("Internal error on token service: %s", err)
		return
	}
	err = pem.Encode(pubKeyFile, &pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyEncoded})
	if err != nil {
		logging.Error().Printf("Internal error on token service: %s", err)
		return
	}

	if viper.GetBool("localNode") {
		daemon.SetPublicKey(pubKey)
	}

	return
}

func generatePrivateKey() (privKey *ecdsa.PrivateKey, err error) {
	var key bytes.Buffer
	privKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return
	}

	privKeyEncoded, _ := x509.MarshalECPrivateKey(privKey)
	privKeyFile, err := os.OpenFile(viper.GetString("panel.token.private"), os.O_CREATE|os.O_WRONLY, 0600)
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
