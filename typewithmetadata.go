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
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

//designed to be overridden
type MetadataType struct {
	Type     string                 `json:"type,omitempty"`
	Metadata map[string]interface{} `json:"-,omitempty"`
}

//parses a type with this declaration, storing what it needs into metadata and type
func (t *MetadataType) UnmarshalJSON(bs []byte) (err error) {
	err = json.Unmarshal(bs, &t.Metadata)
	if err != nil {
		return
	}
	a := t.Metadata["type"]
	if a == nil {
		return errors.New("no type defined")
	}
	var ok bool
	t.Type, ok = a.(string)
	if !ok {
		return errors.New(fmt.Sprintf("type is of %s instead of string", reflect.TypeOf(a)))
	}

	delete(t.Metadata, "type")
	return
}

func (t *MetadataType) MarshalJSON() ([]byte, error) {
	newMapping := make(map[string]interface{})
	for k, v := range t.Metadata {
		newMapping[k] = v
	}
	newMapping["type"] = t.Type
	return json.Marshal(newMapping)
}

//Parses the metadata into the target interface
func (t *MetadataType) ParseMetadata(target interface{}) (err error) {
	data, err := json.Marshal(t)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &target)
	return
}
