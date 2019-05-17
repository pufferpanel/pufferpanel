package view

import (
	"github.com/pufferpanel/pufferpanel/errors"
	"github.com/pufferpanel/pufferpanel/models"
	"gopkg.in/go-playground/validator.v9"
	"net/url"
)

type UserViewModel struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	//ONLY SHOW WHEN COPYING
	Password string `json:"password,omitempty"`
}

func FromUser(model *models.User) *UserViewModel {
	return &UserViewModel{
		Username: model.Username,
		Email:    model.Email,
	}
}

func FromUsers(users *models.Users) []*UserViewModel {
	result := make([]*UserViewModel, len(*users))

	for k, v := range *users {
		result[k] = FromUser(v)
	}

	return result
}

func (model *UserViewModel) CopyToModel(newModel *models.User) {
	if model.Username != "" {
		newModel.Username = model.Username
	}

	if model.Email != "" {
		newModel.Email = model.Email
	}

	if model.Password != "" {
		newModel.SetPassword(model.Password)
	}
}

func (model *UserViewModel) Valid(allowEmpty bool) error {
	validate := validator.New()

	if !allowEmpty && validate.Var(model.Username, "required") != nil {
		return errors.ErrFieldRequired("username")
	}

	if validate.Var(model.Username, "optional|printascii") != nil {
		return errors.ErrFieldMustBePrintable("username")
	}

	testName := url.QueryEscape(model.Username)
	if testName != model.Username {
		return errors.ErrFieldHasURICharacters("username")
	}

	if !allowEmpty && validate.Var(model.Email, "required") != nil {
		return errors.ErrFieldRequired("email")
	}

	if validate.Var(model.Email, "optional|email") != nil {
		return errors.ErrFieldNotEmail("email")
	}

	return nil
}