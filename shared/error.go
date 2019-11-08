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

package shared

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
