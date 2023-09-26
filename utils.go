package pufferpanel

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func Union[T comparable](a, b []T) []T {
	result := make([]T, 0)

	if a == nil || b == nil || len(a) == 0 || len(b) == 0 {
		return result
	}

	for _, v := range a {
		for _, x := range b {
			if v == x {
				result = append(result, v)
				break
			}
		}
	}

	return result
}
