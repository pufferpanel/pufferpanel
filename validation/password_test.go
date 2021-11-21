package validation

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
	"testing"
)

type testPassword struct {
	Password string `validate:"required,entropy"`
}

func TestEntropy(t *testing.T) {
	validate := validator.New()
	assert.NoError(t, validate.RegisterValidation("entropy", PasswordEntropy))
	assert.NoError(t, validate.Struct(testPassword{
		Password: "adsh&ASYd8a78sdhga8",
	}))
	assert.NoError(t, validate.Struct(testPassword{
		Password: "longtextpassword",
	}))
}

func TestEntropyFailure(t *testing.T) {
	validate := validator.New()
	assert.NoError(t, validate.RegisterValidation("entropy", PasswordEntropy))
	assert.Error(t, validate.Struct(testPassword{
		Password: "password",
	}))
	assert.Error(t, validate.Struct(testPassword{
		Password: "PASSWORD",
	}))
	assert.Error(t, validate.Struct(testPassword{
		Password: "PaSsWoRD",
	}))
	assert.Error(t, validate.Struct(testPassword{
		Password: "abc123",
	}))
	assert.Error(t, validate.Struct(testPassword{
		Password: "Password1.",
	}))
}
