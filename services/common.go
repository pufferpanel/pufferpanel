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

package services

import "strings"

func ParseAllowedTags(source string, allowed []string) []string {
	includeTags := make([]string, 0)

	args := strings.Split(strings.ToLower(source), ",")

	for _, test := range args {
		for _, v := range allowed {
			if test == v {
				includeTags = append(includeTags, v)
			}
		}
	}

	return includeTags
}
