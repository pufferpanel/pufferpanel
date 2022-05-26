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

package response

import (
	"github.com/pufferpanel/pufferpanel/v2"
)

type Error struct {
	Error *pufferpanel.Error `json:"error"`
}

type Metadata struct {
	Paging *Paging `json:"paging"`
}

type Paging struct {
	Page    uint  `json:"page"`
	Size    uint  `json:"pageSize"`
	MaxSize uint  `json:"maxSize"`
	Total   int64 `json:"total"`
}

type Empty struct {
}
