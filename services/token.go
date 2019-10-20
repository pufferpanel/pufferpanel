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
	"github.com/pufferpanel/apufferi/v3"
	"github.com/pufferpanel/apufferi/v3/logging"
	"github.com/pufferpanel/apufferi/v3/scope"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/models"
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
	claims := &apufferi.Claim{
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
	claims := &apufferi.Claim{
		StandardClaims: jwt.StandardClaims{
			Audience:  "oauth2",
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		PanelClaims: apufferi.PanelClaims{
			Scopes: map[string][]scope.Scope{
				client.ServerId: client.Scopes,
			},
		},
	}

	return Generate(claims)
}

func GenerateOAuthForNode(nodeId uint) (string, error) {
	claims := &apufferi.Claim{
		StandardClaims: jwt.StandardClaims{
			Audience:  "oauth2",
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		PanelClaims: apufferi.PanelClaims{
			Scopes: map[string][]scope.Scope{
				"": {scope.OAuth2Auth},
			},
		},
	}
	return Generate(claims)
}

func GenerateOAuthForUser(userId uint, serverId *string) (string, error) {
	db, err := database.GetConnection()
	if err != nil {
		return "", err
	}
	ps := &Permission{DB: db}

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

	claims := &apufferi.Claim{
		StandardClaims: jwt.StandardClaims{
			Audience:  "oauth2",
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		PanelClaims: apufferi.PanelClaims{
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

func ParseToken(token string) (*apufferi.Token, error) {
	ValidateTokenLoaded()
	return apufferi.ParseToken(&privateKey.PublicKey, token)
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
	privKeyFile, err := os.OpenFile(viper.GetString("token.private"), os.O_RDONLY, 0600)
	defer apufferi.Close(privKeyFile)
	if os.IsNotExist(err) {
		privKey, err = generatePrivateKey()
	} else if err == nil {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, privKeyFile)
		block, _ := pem.Decode(buf.Bytes())

		privKey, err = ecdsa.GenerateKey(elliptic.P256(), bytes.NewReader(block.Bytes))
	}

	if err != nil {
		logging.Build(logging.ERROR).WithMessage("internal error on token service").WithError(err).Log()
		return
	}

	privateKey = privKey

	pubKey := &privateKey.PublicKey
	pubKeyEncoded, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		logging.Build(logging.ERROR).WithMessage("internal error on token service").WithError(err).Log()
		return
	}

	pubKeyFile, err := os.OpenFile(viper.GetString("token.public"), os.O_CREATE|os.O_RDWR, 0644)
	defer apufferi.Close(pubKeyFile)
	if err != nil {
		logging.Build(logging.ERROR).WithMessage("internal error on token service").WithError(err).Log()
		return
	}
	err = pem.Encode(pubKeyFile, &pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyEncoded})
	if err != nil {
		logging.Build(logging.ERROR).WithMessage("internal error on token service").WithError(err).Log()
		return
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
	privKeyFile, err := os.OpenFile(viper.GetString("token.private"), os.O_CREATE|os.O_WRONLY, 0600)
	defer apufferi.Close(privKeyFile)
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
