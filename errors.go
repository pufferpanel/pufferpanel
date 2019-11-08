package pufferpanel

import (
	"errors"
	"github.com/pufferpanel/pufferpanel/v2/shared"
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

var ErrUnknownError = shared.CreateError("unknown error", "ErrUnknownError")
var ErrInvalidCredentials = shared.CreateError("invalid credentials", "ErrInvalidCredentials")
var ErrServiceNotAvailable = shared.CreateError("service not available", "ErrServiceNotAvailable")
var ErrEmailNotConfigured = shared.CreateError("email not configured", "ErrEmailNotConfigured")
var ErrTokenInvalid = shared.CreateError("token is invalid", "ErrTokenInvalid")
var ErrClientNotFound = shared.CreateError("client not found", "ErrClientNotFound")
var ErrUserNotFound = shared.CreateError("user not found", "ErrUserNotFound")
var ErrLoginNotPermitted = shared.CreateError("login not permitted", "ErrLoginNotPermitted")
var ErrInvalidTokenState = shared.CreateError("invalid token state", "ErrInvalidTokenState")
var ErrNoPermission = shared.CreateError("no permission to perform action", "ErrNoPermission")
var ErrServerNotFound = shared.CreateError("server not found", "ErrServerNotFound")
var ErrDatabaseNotAvailable = shared.CreateError("database not available", "ErrDatabaseNotAvailable")

var ErrSettingNotConfigured = func(name string) *shared.Error {
	return shared.CreateError("${setting} is not configured", "ErrSettingNotConfigured").Metadata(map[string]interface{}{"setting": name})
}

var ErrNoTemplate = func(template string) *shared.Error {
	return shared.CreateError("no template with name ${name}", "ErrNoTemplate").Metadata(map[string]interface{}{"template": template})
}

var ErrServiceInvalidProvider = func(service, provider string) *shared.Error {
	return shared.CreateError("{service} does not support ${provider}", "ErrServiceInvalidProvider").Metadata(map[string]interface{}{"service": service, "provider": provider})
}

var ErrFieldRequired = func(fieldName string) *shared.Error {
	return shared.CreateError("${field} is required", "ErrFieldRequired").Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldMustBePrintable = func(fieldName string) *shared.Error {
	return shared.CreateError("${field} must be printable ascii characters", "ErrFieldMustBePrintable").Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldHasURICharacters = func(fieldName string) *shared.Error {
	return shared.CreateError("${field} must not contain characters which cannot be used in URIs", "ErrFieldHasURICharacters").Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldIsInvalidHost = func(fieldName string) *shared.Error {
	return shared.CreateError("${field} must be a valid IP or FQDN", "ErrFieldIsInvalidHost").Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldIsInvalidIP = func(fieldName string) *shared.Error {
	return shared.CreateError("${field} must be a valid IP", "ErrFieldIsInvalidIP").Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldTooLarge = func(fieldName string, value int64) *shared.Error {
	return shared.CreateError("${field} cannot be larger than ${max}", "ErrFieldTooLarge").Metadata(map[string]interface{}{"field": fieldName, "max": value})
}

var ErrFieldTooSmall = func(fieldName string, value int64) *shared.Error {
	return shared.CreateError("${field} cannot be smaller than ${min}", "ErrFieldTooSmall").Metadata(map[string]interface{}{"field": fieldName, "min": value})
}

var ErrFieldNotBetween = func(fieldName string, min, max int64) *shared.Error {
	return shared.CreateError("${field} must be between ${min} and ${max}", "ErrFieldNotBetween").Metadata(map[string]interface{}{"field": fieldName, "min": min, "max": max})
}

var ErrFieldEqual = func(fieldName1, fieldName2 string) *shared.Error {
	return shared.CreateError("${field1} cannot be equal to ${field2}", "ErrFieldEqual").Metadata(map[string]interface{}{"field1": fieldName1, "field2": fieldName2})
}

var ErrFieldNotEqual = func(fieldName1, fieldName2 string) *shared.Error {
	return shared.CreateError("${field1} is not equal to ${field2}", "ErrFieldNotEqual").Metadata(map[string]interface{}{"field1": fieldName1, "field2": fieldName2})
}

var ErrFieldNotEmail = func(fieldName string) *shared.Error {
	return shared.CreateError("${field} is not a valid email", "ErrFieldNotEmail").Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldLength = func(fieldName string, length int) *shared.Error {
	return shared.CreateError("${field} must be at least ${length} characters", "ErrFieldLength").Metadata(map[string]interface{}{"field": fieldName, "length": length})
}

var ErrNodeInvalid = shared.CreateError("node is invalid", "ErrNodeInvalid")

func GenerateValidationMessage(err error) error {
	if errs, ok := err.(validator.ValidationErrors); ok {
		msg := make([]string, 0)
		for _, e := range errs {
			t := e.Field() + ": " + e.ActualTag()
			if e.Param() != "" {
				t += " (" + e.Param() + ")"
			}
			msg = append(msg, t)
		}
		return errors.New(strings.Join(msg, ", "))
	}
	return nil
}
