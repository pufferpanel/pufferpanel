package services

import (
	"context"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"net/http"
	"sync"
)

type TokenService interface {
	GetKeyFunc() jwt.Keyfunc
	GetTokenStore() jwkset.Storage
	GenerateRequest() (string, error)
	ValidateRequest(string) error
}

type tokenService struct{}

var externalService keyfunc.Keyfunc
var tokenServiceLocker sync.Mutex
var tokenStore jwkset.Storage
var privateKey crypto.PrivateKey

const keyId = "pufferpanel"

func NewTokenService() (TokenService, error) {
	tokenServiceLocker.Lock()
	defer tokenServiceLocker.Unlock()

	if config.PanelEnabled.Value() {
		if tokenStore == nil {
			key := config.PrivateKey.Value()

			var randData []byte
			var err error
			if key == "" {
				randData = make([]byte, ed25519.SeedSize)
				_, err = rand.Read(randData)
				if err != nil {
					return nil, err
				}
				enc := base64.StdEncoding.EncodeToString(randData)
				_ = config.PrivateKey.Set(enc, true)
			} else {
				randData, err = base64.StdEncoding.DecodeString(key)
				if err != nil {
					return nil, err
				}
			}

			// Create a cryptographic key.
			priv := ed25519.NewKeyFromSeed(randData)
			pub := priv.Public()
			privateKey = priv

			// Turn the key into a JWK.
			options := jwkset.JWKOptions{
				Marshal: jwkset.JWKMarshalOptions{
					Private: false,
				},
				Metadata: jwkset.JWKMetadataOptions{
					KID: keyId,
				},
			}

			privateJwk, err := jwkset.NewJWKFromKey(pub, options)
			if err != nil {
				return nil, err
			}

			tokenStore = jwkset.NewMemoryStorage()
			err = tokenStore.KeyWrite(context.Background(), privateJwk)
			if err != nil {
				return nil, err
			}
		}

		if externalService == nil {
			var err error
			externalService, err = keyfunc.New(keyfunc.Options{
				Storage: tokenStore,
			})
			if err != nil {
				return nil, err
			}
		}
	} else {
		if externalService == nil {
			var err error
			externalService, err = keyfunc.NewDefault([]string{})
			if err != nil {
				return nil, err
			}
		}
	}

	return &tokenService{}, nil
}

func TokenServiceGetPublicKey(c *gin.Context) {
	ts, err := NewTokenService()
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}
	rawJWKS, err := ts.GetTokenStore().JSONPublic(context.Background())
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}
	c.JSON(http.StatusOK, rawJWKS)
}

func (ts *tokenService) GenerateRequest() (string, error) {
	token := jwt.New(jwt.SigningMethodEdDSA)
	token.Header[jwkset.HeaderKID] = keyId
	token.Claims = jwt.MapClaims{}

	signed, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (ts *tokenService) ValidateRequest(token string) error {
	parsed, err := jwt.Parse(token, ts.GetKeyFunc())
	if err != nil {
		return err
	}
	if !parsed.Valid {
		return jwt.ErrTokenSignatureInvalid
	}
	return nil
}

func (ts *tokenService) GetKeyFunc() jwt.Keyfunc {
	return externalService.Keyfunc
}

func (ts *tokenService) GetTokenStore() jwkset.Storage {
	return tokenStore
}
