package auth

import (
	"github.com/pufferpanel/pufferpanel/v2"
	"gopkg.in/go-playground/validator.v9"
	"math"
	"reflect"
)

const EntropyThreshold = 32

func PasswordEntropy(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.String {
		return false
	}

	password := fl.Field().String()

	return GetEntropy(password) > EntropyThreshold
}

func EntropyWithErr(password string) error {
	entropy := GetEntropy(password)
	if entropy < EntropyThreshold {
		return pufferpanel.ErrLowEntropy(entropy, EntropyThreshold)
	}
	return nil
}

func GetEntropy(password string) float64 {
	buckets := make(map[rune]bool)
	// Collect all the unique symbols
	for _, char := range password {
		buckets[char] = true
	}

	// log2(R^L)
	// r = symbol quantity
	// l = password length
	return math.Log2(math.Pow(float64(len(buckets)), float64(len(password))))
}
