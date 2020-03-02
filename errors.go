/*
 Copyright 2020 Padduck, LLC
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

package pufferpanel

import (
	"errors"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"gopkg.in/go-playground/validator.v9"
	"runtime/debug"
	"strings"
)

var ErrUnknownError = CreateError("unknown error", "ErrUnknownError")
var ErrInvalidCredentials = CreateError("invalid credentials", "ErrInvalidCredentials")
var ErrServiceNotAvailable = CreateError("service not available", "ErrServiceNotAvailable")
var ErrEmailNotConfigured = CreateError("email not configured", "ErrEmailNotConfigured")
var ErrTokenInvalid = CreateError("token is invalid", "ErrTokenInvalid")
var ErrClientNotFound = CreateError("client not found", "ErrClientNotFound")
var ErrUserNotFound = CreateError("user not found", "ErrUserNotFound")
var ErrLoginNotPermitted = CreateError("login not permitted", "ErrLoginNotPermitted")
var ErrInvalidTokenState = CreateError("invalid token state", "ErrInvalidTokenState")
var ErrNoPermission = CreateError("no permission to perform action", "ErrNoPermission")
var ErrServerNotFound = CreateError("server not found", "ErrServerNotFound")
var ErrDatabaseNotAvailable = CreateError("database not available", "ErrDatabaseNotAvailable")
var ErrServerOffline = CreateError("server offline", "ErrServerOffline")
var ErrIllegalFileAccess = CreateError("invalid file access", "ErrIllegalFileAccess")
var ErrServerDisabled = CreateError("server is disabled", "ErrServerDisabled")
var ErrContainerRunning = CreateError("container already running", "ErrContainerRunning")
var ErrImageDownloading = CreateError("image downloading", "ErrImageDownloading")
var ErrProcessRunning = CreateError("process already running", "ErrProcessRunning")
var ErrMissingFactory = CreateError("missing factory", "ErrMissingFactory")
var ErrServerAlreadyExists = CreateError("server already exists", "ErrServerAlreadyExists")
var ErrInvalidUnixTime = CreateError("time provided is not a valid UNIX time", "ErrInvalidUnixTime")
var ErrKeyNotPEM = CreateError("key is not in PEM format", "ErrKeyNotPEM")
var ErrCannotValidateToken = CreateError("could not validate access token", "ErrCannotValidateToken")
var ErrMissingAccessToken = CreateError("access token not provided", "ErrMissingAccessToken")
var ErrNotBearerToken = CreateError("access token must be a Bearer token", "ErrNotBearerToken")
var ErrKeyNotECDSA = CreateError("key is not ECDSA key", "ErrKeyNotECDSA")
var ErrMissingScope = CreateError("missing scope", "ErrMissingScope")

func CreateErrMissingScope(scope Scope) *Error {
	return CreateError(ErrMissingScope.Message, ErrMissingScope.Code).Metadata(map[string]interface{}{"scope": scope})
}

var ErrSettingNotConfigured = func(name string) *Error {
	return CreateError("${setting} is not configured", "ErrSettingNotConfigured").Metadata(map[string]interface{}{"setting": name})
}

var ErrNoTemplate = func(template string) *Error {
	return CreateError("no template with name ${name}", "ErrNoTemplate").Metadata(map[string]interface{}{"template": template})
}

var ErrServiceInvalidProvider = func(service, provider string) *Error {
	return CreateError("{service} does not support ${provider}", "ErrServiceInvalidProvider").Metadata(map[string]interface{}{"service": service, "provider": provider})
}

var ErrFieldRequired = func(fieldName string) *Error {
	return CreateError("${field} is required", "ErrFieldRequired").Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldMustBePrintable = func(fieldName string) *Error {
	return CreateError("${field} must be printable ascii characters", "ErrFieldMustBePrintable").Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldHasURICharacters = func(fieldName string) *Error {
	return CreateError("${field} must not contain characters which cannot be used in URIs", "ErrFieldHasURICharacters").Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldIsInvalidHost = func(fieldName string) *Error {
	return CreateError("${field} must be a valid IP or FQDN", "ErrFieldIsInvalidHost").Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldIsInvalidIP = func(fieldName string) *Error {
	return CreateError("${field} must be a valid IP", "ErrFieldIsInvalidIP").Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldTooLarge = func(fieldName string, value int64) *Error {
	return CreateError("${field} cannot be larger than ${max}", "ErrFieldTooLarge").Metadata(map[string]interface{}{"field": fieldName, "max": value})
}

var ErrFieldTooSmall = func(fieldName string, value int64) *Error {
	return CreateError("${field} cannot be smaller than ${min}", "ErrFieldTooSmall").Metadata(map[string]interface{}{"field": fieldName, "min": value})
}

var ErrFieldNotBetween = func(fieldName string, min, max int64) *Error {
	return CreateError("${field} must be between ${min} and ${max}", "ErrFieldNotBetween").Metadata(map[string]interface{}{"field": fieldName, "min": min, "max": max})
}

var ErrFieldEqual = func(fieldName1, fieldName2 string) *Error {
	return CreateError("${field1} cannot be equal to ${field2}", "ErrFieldEqual").Metadata(map[string]interface{}{"field1": fieldName1, "field2": fieldName2})
}

var ErrFieldNotEqual = func(fieldName1, fieldName2 string) *Error {
	return CreateError("${field1} is not equal to ${field2}", "ErrFieldNotEqual").Metadata(map[string]interface{}{"field1": fieldName1, "field2": fieldName2})
}

var ErrFieldNotEmail = func(fieldName string) *Error {
	return CreateError("${field} is not a valid email", "ErrFieldNotEmail").Metadata(map[string]interface{}{"field": fieldName})
}

var ErrFieldLength = func(fieldName string, min int, max int) *Error {
	return CreateError("${field} must be between ${min} and ${max} characters", "ErrFieldLength").Metadata(map[string]interface{}{"field": fieldName, "min": min, "max": max})
}

var ErrNodeInvalid = CreateError("node is invalid", "ErrNodeInvalid")

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

type Error struct {
	Message string                 `json:"msg,omitempty"`
	Code    string                 `json:"code,omitempty"`
	Meta    map[string]interface{} `json:"metadata,omitempty"`
	error
}

func (ge *Error) GetMessage() string {
	return ReplaceTokens(ge.Message, ge.Meta)
}

func (ge *Error) GetCode() string {
	return ge.Code
}

func (ge *Error) Error() string {
	return ge.GetMessage()
}

func (ge *Error) Is(err *Error) bool {
	return ge.GetCode() == err.GetCode()
}

func (ge *Error) Metadata(metadata map[string]interface{}) *Error {
	cp := ge
	cp.Meta = metadata
	return cp
}

func CreateError(msg, code string) *Error {
	return &Error{
		Message: msg,
		Code:    code,
	}
}

func FromError(err error) *Error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*Error); ok {
		return e
	}
	return CreateError(err.Error(), "ErrGeneric")
}

func Recover() {
	if err := recover(); err != nil {
		if _, ok := err.(error); !ok {
			err = errors.New(ToString(err))
		}

		logging.Error().Printf("CRITICAL ERROR: \n%+v\n%s", err, debug.Stack())
	}
}
