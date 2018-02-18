package models

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/satori/go.uuid"
)

func TestLocation_ValidatePass(t *testing.T) {
	uuid, err := uuid.NewV4()
	loc := &Location{
		ID: uuid,
		Code: "test",
		Name: "Code Test",
	}

	validationErrors, err := loc.Validate(nil)
	assert.NoError(t, err, "location is not valid when it should be")
	assert.Falsef(t, validationErrors.HasAny(), "location is not valid: (%v)", validationErrors)
}

func TestLocation_ValidateFailNoCode(t *testing.T) {
	uuid, _ := uuid.NewV4()
	loc := &Location{
		ID: uuid,
		Name: "Code Test",
	}

	validationErrors, err := loc.Validate(nil)
	assert.NoError(t, err, "error thrown when there should have been", err)
	assert.True(t, validationErrors.HasAny(), "errors were not returned from validation", validationErrors)
}

func TestLocation_ValidateFailNoName(t *testing.T) {
	uuid, err := uuid.NewV4()
	loc := &Location{
		ID: uuid,
		Code: "test",
	}

	validationErrors, err := loc.Validate(nil)
	assert.NoError(t, err, "error thrown when there should have been", err)
	assert.True(t, validationErrors.HasAny(), "errors were not returned from validation", validationErrors)
}

func TestLocation_ValidateFailNoId(t *testing.T) {
	loc := &Location{
		Code: "test",
		Name: "Code Test",
	}

	validationErrors, err := loc.Validate(nil)
	assert.NoError(t, err, "error thrown when there should have been", err)
	assert.True(t, validationErrors.HasAny(), "errors were not returned from validation", validationErrors)
}