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

package models

import (
	"encoding/json"
	"github.com/pufferpanel/apufferi/v4"
	"strings"
)

//this is basically the template, just wrapped enough to be used in the database
type Template struct {
	apufferi.Template `gorm:"-"`

	ID       uint   `json:"-"`
	Name     string `gorm:"type:varchar(100)" json:"name"`
	RawValue string `gorm:"type:text" json:"-"`

	Readme string `gorm:"type:text" json:"readme,omitempty"`
}

type Templates []*Template

func (t *Template) AfterFind() error {
	err := json.NewDecoder(strings.NewReader(t.RawValue)).Decode(&t.Template)
	if err != nil {
		return err
	}
	t.RawValue = ""
	return nil
}

func (t *Template) BeforeSave() error {
	data, err := json.Marshal(&t.Template)
	if err != nil {
		return err
	}
	t.RawValue = string(data)
	return nil
}

func (ts *Templates) MarshalJSON() ([]byte, error) {
	res := make(map[string]*Template)

	for _, v := range *ts {
		res[v.Name] = v
	}

	return json.Marshal(res)
}
