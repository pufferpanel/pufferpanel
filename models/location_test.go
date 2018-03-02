package models

import (
	"github.com/gobuffalo/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLocation_ValidatePass(t *testing.T) {
	uuid, err := uuid.NewV4()
	loc := &Location{
		ID:   uuid,
		Code: "test",
		Name: "Code Test",
	}

	validationErrors, err := loc.Validate(nil)
	assert.NoError(t, err, "error thrown when there should have been", err)
	assert.Falsef(t, validationErrors.HasAny(), "location is not valid: (%v)", validationErrors)
}

func TestLocation_ValidateFailNoCode(t *testing.T) {
	uuid, _ := uuid.NewV4()
	loc := &Location{
		ID:   uuid,
		Name: "Code Test",
	}

	validationErrors, err := loc.Validate(nil)
	assert.NoError(t, err, "error thrown when there should have been", err)
	assert.True(t, validationErrors.HasAny(), "errors were not returned from validation", validationErrors)
}

func TestLocation_ValidateFailNoName(t *testing.T) {
	uuid, err := uuid.NewV4()
	loc := &Location{
		ID:   uuid,
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
