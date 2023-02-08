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

package spongedl

import (
	"encoding/json"
	"github.com/pufferpanel/pufferpanel/v3"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Key() string {
	return "spongedl"
}
func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	jsonData, err := json.Marshal(op.OperationArgs)
	if err != nil {
		return nil, err
	}

	var spongedl SpongeDl
	err = json.Unmarshal(jsonData, &spongedl)
	return spongedl, err
}

var Factory OperationFactory
