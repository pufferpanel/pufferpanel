package auth

import (
	"gopkg.in/go-playground/validator.v9"
	"math"
	"reflect"
)

func passwordEntropy(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.String {
		return false
	}

	password := fl.Field().String()
	buckets := make(map[rune]bool)
	// Collect all the unique symbols
	for _, char := range password {
		buckets[char] = true
	}

	// log2(R^L)
	// r = symbol quantity
	// l = password length
	entropy := math.Log2(math.Pow(float64(len(buckets)), float64(len(password))))
	return entropy > 32
}
