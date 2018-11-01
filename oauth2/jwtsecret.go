package oauth2

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"strings"
)

var secretToken string

func GetJWTSecret() string {
	if secretToken == "" {
		GenerateSecret()
	}
	return secretToken
}

func GenerateSecret() {
	secretToken = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	fmt.Printf("Generated secure token for JWT: %s\n", secretToken)
}

func Validate(token string) bool {
	jwt := NewJWTAccessGenerate([]byte(GetJWTSecret()), jwt.SigningMethodHS512)
	return jwt.Validate(token)
}