/*
 Copyright 2019 Padduck, LLC
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at
  	http://www.apache.org/licenses/LICENSE-2.0
  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package errors

import "github.com/pufferpanel/apufferi"

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

var ErrSettingNotConfigured = func(name string) apufferi.Error {
	return apufferi.CreateError("%v is not configured", "ErrSettingNotConfigured").Set(name).Metadata(map[string]interface{}{"setting": name})
}

var ErrNoTemplate = func(template string) apufferi.Error {
	return apufferi.CreateError("no template with given name", "ErrNoTemplate").Metadata(map[string]interface{}{"template": template})
}

var ErrServiceInvalidProvider = func(service, provider string) apufferi.Error {
	return apufferi.CreateError("%v does not support %v", "ErrServiceInvalidProvider").Set(service, provider).Metadata(map[string]interface{}{"service": service, "provider": provider})
}

var ErrFieldRequired = func(fieldName string) apufferi.Error {
	return apufferi.CreateError("%v is required", "ErrFieldRequired").Set(fieldName).Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldMustBePrintable = func(fieldName string) apufferi.Error {
	return apufferi.CreateError("%v must be printable ascii characters", "ErrFieldMustBePrintable").Set(fieldName).Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldHasURICharacters = func(fieldName string) apufferi.Error {
	return apufferi.CreateError("%v must not contain characters which cannot be used in URIs", "ErrFieldHasURICharacters").Set(fieldName).Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldIsInvalidHost = func(fieldName string) apufferi.Error {
	return apufferi.CreateError("%v must be a valid IP or FQDN", "ErrFieldIsInvalidHost").Set(fieldName).Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldIsInvalidIP = func(fieldName string) apufferi.Error {
	return apufferi.CreateError("%v must be a valid IP", "ErrFieldIsInvalidIP").Set(fieldName).Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldTooLarge = func(fieldName string, value int64) apufferi.Error {
	return apufferi.CreateError("%v cannot be larger than %v", "ErrFieldTooLarge").Set(fieldName, value).Metadata(map[string]interface{}{"field": fieldName, "max": value})
}

var ErrFieldTooSmall = func(fieldName string, value int64) apufferi.Error {
	return apufferi.CreateError("%v cannot be smaller than %v", "ErrFieldTooSmall").Set(fieldName, value).Metadata(map[string]interface{}{"field": fieldName, "min": value})
}

var ErrFieldNotBetween = func(fieldName string, min, max int64) apufferi.Error {
	return apufferi.CreateError("%v must be between %v and %v", "ErrFieldNotBetween").Set(fieldName, min, max).Metadata(map[string]interface{}{"field": fieldName, "min": min, "max": max})
}

var ErrFieldEqual = func(fieldName1, fieldName2 string) apufferi.Error {
	return apufferi.CreateError("%v cannot be equal to %v", "ErrFieldEqual").Set(fieldName1, fieldName2).Metadata(map[string]interface{}{"field1": fieldName1, "field2": fieldName2})
}

var ErrFieldNotEqual = func(fieldName1, fieldName2 string) apufferi.Error {
	return apufferi.CreateError("%v is not equal to %v", "ErrFieldNotEqual").Set(fieldName1, fieldName2).Metadata(map[string]interface{}{"field1": fieldName1, "field2": fieldName2})
}

var ErrFieldNotEmail = func(fieldName string) apufferi.Error {
	return apufferi.CreateError("%v is not a valid email", "ErrFieldNotEmail").Set(fieldName).Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldLength = func(fieldName string, length int) apufferi.Error {
	return apufferi.CreateError("%v must be at least %v characters", "ErrFieldLength").Set(fieldName).Metadata(map[string]interface{}{"field": fieldName, "length": length})
}