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
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func Union[T comparable](a, b []T) []T {
	result := make([]T, 0)

	if a == nil || b == nil || len(a) == 0 || len(b) == 0 {
		return result
	}

	for _, v := range a {
		for _, x := range b {
			if v == x {
				result = append(result, v)
				break
			}
		}
	}

	return result
}
