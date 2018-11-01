package oauth2

import (
	"fmt"
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