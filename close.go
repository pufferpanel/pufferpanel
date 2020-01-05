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

package pufferpanel

import (
	"io"
	"net/http"
)

func Close(closer io.Closer) {
	//at this point, i give up trying to get this to not fail, so we'll go with the brute force
	defer func() {
		recover()
	}()

	if closer != nil {
		_ = closer.Close()
	}

	/*if closer != nil && !(reflect.ValueOf(closer).Kind() != reflect.Ptr && reflect.ValueOf(closer).IsNil()) {
		_ = closer.Close()
	}*/
}

func CloseResponse(response *http.Response) {
	if response != nil {
		Close(response.Body)
	}
}
