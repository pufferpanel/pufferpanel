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
	"github.com/pufferpanel/pufferpanel/v2"
	"gorm.io/gorm"
	"strings"
)

type Template struct {
	pufferpanel.Server `gorm:"-"`

	Name     string `gorm:"type:varchar(100);primaryKey" json:"name"`
	RawValue string `gorm:"type:text" json:"-"`

	Readme string `gorm:"type:text" json:"readme,omitempty"`
}

type Templates []*Template

func (t *Template) AfterFind(*gorm.DB) error {
	err := json.NewDecoder(strings.NewReader(t.RawValue)).Decode(&t.Server)
	if err != nil {
		return err
	}
	t.RawValue = ""
	if t.Execution.LegacyRun != "" {
		t.Execution.Command = strings.TrimSpace(t.Execution.LegacyRun + " " + strings.Join(t.Execution.LegacyArguments, " "))
		t.Execution.LegacyRun = ""
		t.Execution.LegacyArguments = nil
	}
	return nil
}

func (t *Template) BeforeSave(*gorm.DB) error {
	if t.Execution.LegacyRun != "" {
		t.Execution.Command = strings.TrimSpace(t.Execution.LegacyRun + " " + strings.Join(t.Execution.LegacyArguments, " "))
		t.Execution.LegacyRun = ""
		t.Execution.LegacyArguments = nil
	}
	data, err := json.Marshal(&t.Server)
	if err != nil {
		return err
	}
	t.RawValue = string(data)
	return nil
}
