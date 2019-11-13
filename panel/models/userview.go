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
	"github.com/pufferpanel/pufferpanel/v2"
	"gopkg.in/go-playground/validator.v9"
	"net/url"
)

type UserView struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	//ONLY SHOW WHEN COPYING
	Password    string `json:"password,omitempty"`
	NewPassword string `json:"newPassword,omitempty"`
}

func FromUser(model *User) *UserView {
	return &UserView{
		Username: model.Username,
		Email:    model.Email,
	}
}

func FromUsers(users *Users) []*UserView {
	result := make([]*UserView, len(*users))

	for k, v := range *users {
		result[k] = FromUser(v)
	}

	return result
}

func (model *UserView) CopyToModel(newModel *User) {
	if model.Username != "" {
		newModel.Username = model.Username
	}

	if model.Email != "" {
		newModel.Email = model.Email
	}

	if model.Password != "" {
		_ = newModel.SetPassword(model.Password)
	}
}

func (model *UserView) Valid(allowEmpty bool) error {
	validate := validator.New()

	if !allowEmpty && validate.Var(model.Username, "required") != nil {
		return pufferpanel.ErrFieldRequired("username")
	}

	if validate.Var(model.Username, "omitempty,printascii,min=5,max=100") != nil {
		return pufferpanel.ErrFieldMustBePrintable("username")
	}

	testName := url.QueryEscape(model.Username)
	if testName != model.Username {
		return pufferpanel.ErrFieldHasURICharacters("username")
	}

	if !allowEmpty && validate.Var(model.Email, "required") != nil {
		return pufferpanel.ErrFieldRequired("email")
	}

	if validate.Var(model.Email, "omitempty,email,max=255") != nil {
		return pufferpanel.ErrFieldNotEmail("email")
	}

	return nil
}
