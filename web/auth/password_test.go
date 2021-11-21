package auth

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
	"testing"
)

func TestEntropy(t *testing.T) {
	validate := validator.New()
	assert.NoError(t, validate.RegisterValidation("entropy", PasswordEntropy))
	assert.NoError(t, validate.Struct(registerRequestData{
		Username: "adsdfsdfsdf",
		Email:    "a@b.co",
		Password: "adsh&ASYd8a78sdhga8",
	}))
	assert.NoError(t, validate.Struct(registerRequestData{
		Username: "adsdfsdfsdf",
		Email:    "a@b.co",
		Password: "longtextpassword",
	}))
}

func TestEntropyFailure(t *testing.T) {
	validate := validator.New()
	assert.NoError(t, validate.RegisterValidation("entropy", PasswordEntropy))
	assert.Error(t, validate.Struct(registerRequestData{
		Username: "adsdfsdfsdf",
		Email:    "a@b.co",
		Password: "password",
	}))
	assert.Error(t, validate.Struct(registerRequestData{
		Username: "adsdfsdfsdf",
		Email:    "a@b.co",
		Password: "PASSWORD",
	}))
	assert.Error(t, validate.Struct(registerRequestData{
		Username: "adsdfsdfsdf",
		Email:    "a@b.co",
		Password: "PaSsWoRD",
	}))
	assert.Error(t, validate.Struct(registerRequestData{
		Username: "adsdfsdfsdf",
		Email:    "a@b.co",
		Password: "abc123",
	}))
	assert.Error(t, validate.Struct(registerRequestData{
		Username: "adsdfsdfsdf",
		Email:    "a@b.co",
		Password: "Password1.",
	}))
}
