/*
 Copyright 2018 Padduck, LLC

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

package daemon

import (
	"github.com/pufferpanel/pufferpanel/v2/shared"
	"github.com/pufferpanel/pufferpanel/v2/shared/scope"
)

var ErrServerOffline = shared.CreateError("server offline", "ErrServerOffline")
var ErrIllegalFileAccess = shared.CreateError("invalid file access", "ErrIllegalFileAccess")
var ErrServerDisabled = shared.CreateError("server is disabled", "ErrServerDisabled")
var ErrContainerRunning = shared.CreateError("container already running", "ErrContainerRunning")
var ErrImageDownloading = shared.CreateError("image downloading", "ErrImageDownloading")
var ErrProcessRunning = shared.CreateError("process already running", "ErrProcessRunning")
var ErrMissingFactory = shared.CreateError("missing factory", "ErrMissingFactory")
var ErrServerAlreadyExists = shared.CreateError("server already exists", "ErrServerAlreadyExists")
var ErrInvalidUnixTime = shared.CreateError("time provided is not a valid UNIX time", "ErrInvalidUnixTime")
var ErrKeyNotPEM = shared.CreateError("key is not in PEM format", "ErrKeyNotPEM")
var ErrCannotValidateToken = shared.CreateError("could not validate access token", "ErrCannotValidateToken")
var ErrMissingAccessToken = shared.CreateError("access token not provided", "ErrMissingAccessToken")
var ErrNotBearerToken = shared.CreateError("access token must be a Bearer token", "ErrNotBearerToken")
var ErrKeyNotECDSA = shared.CreateError("key is not ECDSA key", "ErrKeyNotECDSA")
var ErrMissingScope = shared.CreateError("missing scope", "ErrMissingScope")

func CreateErrMissingScope(scope scope.Scope) *shared.Error {
	return shared.CreateError(ErrMissingScope.Message, ErrMissingScope.Code).Metadata(map[string]interface{}{"scope": scope})
}
