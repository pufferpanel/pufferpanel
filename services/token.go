package services

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/dgrijalva/jwt-go"
	"github.com/pufferpanel/apufferi/v3"
	"github.com/pufferpanel/apufferi/v3/logging"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/utils/uuid"
	"os"
	"strings"
	"sync"
	"time"
)

var signingMethod = jwt.SigningMethodES256
var privateKey *ecdsa.PrivateKey
var locker sync.Mutex

func NewJWTAccessGenerate() oauth2.AccessGenerate {
	validateTokenLoaded()
	return &jwtAccessGenerate{}
}

func Generate(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(signingMethod, claims)
	return token.SignedString(privateKey)
}

func ParseToken(token string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return &privateKey.PublicKey, nil
	})
}

func validateTokenLoaded() {
	locker.Lock()
	defer locker.Unlock()
	if privateKey == nil {
		load()
	}
}

func load() {
	var privKey *ecdsa.PrivateKey
	privKeyFile, err := os.OpenFile("private.pem", os.O_CREATE|os.O_RDWR, 0600)
	if os.IsNotExist(err) {
		privKey, err = generatePrivateKey()
	} else if err == nil {
		privKey, err = ecdsa.GenerateKey(elliptic.P256(), privKeyFile)
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

	pubKeyFile, err := os.OpenFile("public.pem", os.O_CREATE|os.O_RDWR, 0644)
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
	privKeyFile, err := os.OpenFile("private.pem", os.O_CREATE|os.O_WRONLY, 0600)
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

type JWTAccessClaims struct {
	jwt.StandardClaims
}

// Valid claims verification
func (a *JWTAccessClaims) Valid() error {
	if time.Unix(a.ExpiresAt, 0).Before(time.Now()) {
		return errors.ErrInvalidAccessToken
	}
	return nil
}

// JWTAccessGenerate generate the jwt access token
type jwtAccessGenerate struct {
}

func (a *jwtAccessGenerate) Token(data *oauth2.GenerateBasic, isGenRefresh bool) (access, refresh string, err error) {
	claims := &JWTAccessClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  data.Client.GetID(),
			Subject:   data.UserID,
			ExpiresAt: data.TokenInfo.GetAccessCreateAt().Add(data.TokenInfo.GetAccessExpiresIn()).Unix(),
		},
	}

	access, err = Generate(claims)
	if err != nil {
		return
	}

	if isGenRefresh {
		refresh = base64.URLEncoding.EncodeToString(uuid.NewSHA1(uuid.Must(uuid.NewRandom()), []byte(access)).Bytes())
		refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	}

	return
}
