package models

import (
	"time"
	"github.com/markbates/pop"
	"errors"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/markbates/validate"
)

type Location struct {
	ID          int       `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Name        string    `json:"name" db:"name"`
	Code        string    `json:"code" db:"code"`
	Description string    `json:"description" db:"description"`
}

type Locations []Location

func (l *Location) Validate(tx *pop.Connection) (*validate.Errors, error) {
	resultErrs := validate.NewErrors()

	err := validation.ValidateStruct(l,
		validation.Field(&l.Name, validation.Required),
		validation.Field(&l.Code, validation.Required, validation.Length(1, 10)),
		validation.Field(&l.Description),
	)

	errs, ok := err.(validation.Errors)

	if ok && (err != nil && errs.Filter() != nil) {
		for k, v := range errs {
			resultErrs.Add(k, v.Error())
		}
	}

	return resultErrs, nil
}

func (l *Location) BeforeDestroy(tx *pop.Connection) error {
	node := &Node{}

	exists, err := tx.BelongsTo(l).Exists(&node)

	if err != nil {
		return err
	} else if exists {
		return errors.New("location is associated with nodes")
	}
	return nil
}
