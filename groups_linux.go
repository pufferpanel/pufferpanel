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
	"os/user"
)

func UserInGroup() bool {
	u, err := user.Current()
	if err != nil {
		return false
	}

	groups, err := u.GroupIds()
	if err != nil {
		return false
	}

	allowedIds := make([]string, 0)

	if expectedGroup, err := user.LookupGroup("pufferpanel"); err == nil {
		allowedIds = append(allowedIds, expectedGroup.Gid)
	}

	if rootGroup, err := user.LookupGroup("root"); err == nil {
		allowedIds = append(allowedIds, rootGroup.Gid)
	}

	for _, v := range groups {
		for _, t := range allowedIds {
			if v == t {
				return true
			}
		}
	}

	return false
}
