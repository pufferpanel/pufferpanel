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

package alterfile

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/spf13/cast"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	file := cast.ToString(op.OperationArgs["file"])
	search := cast.ToString(op.OperationArgs["search"])
	replace := cast.ToString(op.OperationArgs["replace"])
	regex := cast.ToBool(op.OperationArgs["regex"])
	return AlterFile{TargetFile: file, Search: search, Replace: replace, Regex: regex}, nil
}

func (of OperationFactory) Key() string {
	return "alterfile"
}

var Factory OperationFactory
