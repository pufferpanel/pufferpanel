package models

import (
	"github.com/gobuffalo/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

const RawPw = "test"
const HashedPw = "$2a$10$iwgb1AVO0if235/tbdd2H.yx.7DfzT/FfHijXkAL.p2H.1YpuOrcG" //test

func TestUser_ValidatePass(t *testing.T) {
	uuid, err := uuid.NewV4()
	user := &User{
		ID:             uuid,
		Username:       "testuser",
		Email:          "test@example.com",
		HashedPassword: HashedPw,
	}

	validationErrors, err := user.Validate(nil)
	assert.NoError(t, err, "error thrown when there should have been", err)
	assert.Falsef(t, validationErrors.HasAny(), "user is not valid: (%v)", validationErrors)
}

func TestUser_ValidateNoUsername(t *testing.T) {
	uuid, err := uuid.NewV4()
	user := &User{
		ID:             uuid,
		Email:          "test@example.com",
		HashedPassword: HashedPw,
	}

	validationErrors, err := user.Validate(nil)
	assert.NoError(t, err, "error thrown when there should have been", err)
	assert.True(t, validationErrors.HasAny(), "errors were not returned from validation")
}

func TestUser_ValidateNoEmail(t *testing.T) {
	uuid, err := uuid.NewV4()
	user := &User{
		ID:             uuid,
		Username:       "testuser",
		HashedPassword: HashedPw,
	}

	validationErrors, err := user.Validate(nil)
	assert.NoError(t, err, "error thrown when there should have been", err)
	assert.True(t, validationErrors.HasAny(), "errors were not returned from validation")
}

func TestUser_ValidateNoPassword(t *testing.T) {
	uuid, err := uuid.NewV4()
	user := &User{
		ID:       uuid,
		Username: "testuser",
		Email:    "test@example.com",
	}

	validationErrors, err := user.Validate(nil)
	assert.NoError(t, err, "error thrown when there should have been", err)
	assert.True(t, validationErrors.HasAny(), "errors were not returned from validation")
}

func TestUser_ValidateBadEmail(t *testing.T) {
	uuid, err := uuid.NewV4()
	user := &User{
		ID:             uuid,
		Username:       "testuser",
		Email:          "testexample.com",
		HashedPassword: HashedPw,
	}

	validationErrors, err := user.Validate(nil)
	assert.NoError(t, err, "error thrown when there should have been", err)
	assert.True(t, validationErrors.HasAny(), "errors were not returned from validation")
}

func TestUser_ValidatePassChange(t *testing.T) {
	user := &User{}

	err := user.ChangePassword("test")
	assert.NoError(t, err, "error thrown generating new password", err)
}

func TestUser_ValidatePassCheck(t *testing.T) {
	user := &User{
		HashedPassword: HashedPw,
	}

	equal := user.ValidatePassword(RawPw)
	assert.True(t, equal, "hash and raw not equal")

	equal = user.ValidatePassword("wrong")
	assert.False(t, equal, "hash and raw equal")
}
