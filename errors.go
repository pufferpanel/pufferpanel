package pufferpanel

import (
	"errors"
	"github.com/pufferpanel/apufferi/v3"
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

var ErrUnknownError = apufferi.CreateError("unknown error", "ErrUnknownError")
var ErrInvalidCredentials = apufferi.CreateError("invalid credentials", "ErrInvalidCredentials")
var ErrServiceNotAvailable = apufferi.CreateError("service not available", "ErrServiceNotAvailable")
var ErrEmailNotConfigured = apufferi.CreateError("email not configured", "ErrEmailNotConfigured")
var ErrTokenInvalid = apufferi.CreateError("token is invalid", "ErrTokenInvalid")
var ErrClientNotFound = apufferi.CreateError("client not found", "ErrClientNotFound")
var ErrUserNotFound = apufferi.CreateError("user not found", "ErrUserNotFound")
var ErrLoginNotPermitted = apufferi.CreateError("login not permitted", "ErrLoginNotPermitted")
var ErrInvalidTokenState = apufferi.CreateError("invalid token state", "ErrInvalidTokenState")
var ErrNoPermission = apufferi.CreateError("no permission to perform action", "ErrNoPermission")
var ErrServerNotFound = apufferi.CreateError("server not found", "ErrServerNotFound")
var ErrDatabaseNotAvailable = apufferi.CreateError("database not available", "ErrDatabaseNotAvailable")

var ErrSettingNotConfigured = func(name string) apufferi.Error {
	return apufferi.CreateError("{setting} is not configured", "ErrSettingNotConfigured").Metadata(map[string]interface{}{"setting": name})
}

var ErrNoTemplate = func(template string) apufferi.Error {
	return apufferi.CreateError("no template with name {name}", "ErrNoTemplate").Metadata(map[string]interface{}{"template": template})
}

var ErrServiceInvalidProvider = func(service, provider string) apufferi.Error {
	return apufferi.CreateError("{service} does not support {provider}", "ErrServiceInvalidProvider").Metadata(map[string]interface{}{"service": service, "provider": provider})
}

var ErrFieldRequired = func(fieldName string) apufferi.Error {
	return apufferi.CreateError("{field} is required", "ErrFieldRequired").Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldMustBePrintable = func(fieldName string) apufferi.Error {
	return apufferi.CreateError("{field} must be printable ascii characters", "ErrFieldMustBePrintable").Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldHasURICharacters = func(fieldName string) apufferi.Error {
	return apufferi.CreateError("{field} must not contain characters which cannot be used in URIs", "ErrFieldHasURICharacters").Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldIsInvalidHost = func(fieldName string) apufferi.Error {
	return apufferi.CreateError("{field} must be a valid IP or FQDN", "ErrFieldIsInvalidHost").Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldIsInvalidIP = func(fieldName string) apufferi.Error {
	return apufferi.CreateError("{field} must be a valid IP", "ErrFieldIsInvalidIP").Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldTooLarge = func(fieldName string, value int64) apufferi.Error {
	return apufferi.CreateError("{field} cannot be larger than {max}", "ErrFieldTooLarge").Metadata(map[string]interface{}{"field": fieldName, "max": value})
}

var ErrFieldTooSmall = func(fieldName string, value int64) apufferi.Error {
	return apufferi.CreateError("{field} cannot be smaller than {min}", "ErrFieldTooSmall").Metadata(map[string]interface{}{"field": fieldName, "min": value})
}

var ErrFieldNotBetween = func(fieldName string, min, max int64) apufferi.Error {
	return apufferi.CreateError("{field} must be between {min} and {max}", "ErrFieldNotBetween").Metadata(map[string]interface{}{"field": fieldName, "min": min, "max": max})
}

var ErrFieldEqual = func(fieldName1, fieldName2 string) apufferi.Error {
	return apufferi.CreateError("{field1} cannot be equal to {field2}", "ErrFieldEqual").Metadata(map[string]interface{}{"field1": fieldName1, "field2": fieldName2})
}

var ErrFieldNotEqual = func(fieldName1, fieldName2 string) apufferi.Error {
	return apufferi.CreateError("{field1} is not equal to {field2}", "ErrFieldNotEqual").Metadata(map[string]interface{}{"field1": fieldName1, "field2": fieldName2})
}

var ErrFieldNotEmail = func(fieldName string) apufferi.Error {
	return apufferi.CreateError("{field} is not a valid email", "ErrFieldNotEmail").Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldLength = func(fieldName string, length int) apufferi.Error {
	return apufferi.CreateError("{field} must be at least {length} characters", "ErrFieldLength").Metadata(map[string]interface{}{"field": fieldName, "length": length})
}

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
