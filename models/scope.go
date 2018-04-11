package models

import (
	"github.com/gobuffalo/uuid"
	"time"
	"github.com/gobuffalo/pop"
	"github.com/go-ozzo/ozzo-validation"
	"fmt"
	"github.com/gobuffalo/validate"
	"errors"
)

type Scope struct {
	ID        int       `json:"-" db:"id"`
	ClientId  uuid.UUID `json:"-" db:"client_id" rw:"r"`
	ServerId  uuid.UUID `json:"server_id" db:"server_id" rw:"r"`
	value     string    `json:"value" db:"value" rw:"r"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Scopes []Scope

func (s *Scope) SetValue(str string) {
	s.value = str
}

func (s *Scope) GetValue() string {
	return s.ServerId.String() + ":" + s.value
}

func (s *Scope) GetRawValue() string {
	return s.value
}

func (s *Scope) Validate(tx *pop.Connection) (*validate.Errors, error) {
	validationErrors := validate.NewErrors()

	err := validation.ValidateStruct(s,
		validation.Field(&s.ClientId, validation.Required),
		validation.Field(&s.value, validation.Required),
	)
	errs, ok := err.(validation.Errors)

	if err == nil {
		ok = true
	}

	if ok && (err != nil && errs.Filter() != nil) {
		for k, v := range errs {
			validationErrors.Add(k, v.Error())
		}
	} else if !ok {
		return validationErrors, errors.New(fmt.Sprintf("could not cast to validation.Errors (%T)", err))
	}

	return validationErrors, nil
}