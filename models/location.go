package models

import (
	"errors"
	"time"

	"fmt"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/uuid"
)

type Location struct {
	ID        uuid.UUID	`json:"id" db:"id"`
	Code      string    `json:"code" db:"code"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

type Locations []Location

func GetLocations() (locations Locations, err error) {
	locations = Locations{}
	err = DB.All(&locations)
	return
}

func GetLocationById(id string) (location Location, err error) {
	location = Location{}
	err = DB.Find(&location, id)
	return
}

func GetLocationByCode(code string) (location Location, err error) {
	location = Location{}
	query := DB.Where("code = ?", code)
	exists, err := query.Exists(&location)
	if exists {
		err = DB.Where("code = ?", code).First(&location)
	}
	return
}

func CreateLocation(code, name string) (location Location, err error) {
	location = Location{
		Code: code,
		Name: name,
	}

	return
}

func (l *Location) Delete() (err error) {
	err = DB.Destroy(l)
	return
}

func (l *Location) Save() (err error) {
	validationErrors, err := DB.ValidateAndSave(l)
	if validationErrors != nil && validationErrors.Count() > 0 {
		err = errors.New("model is invalid: " + validationErrors.Error())
	}
	return
}

func (l *Location) Validate(tx *pop.Connection) (*validate.Errors, error) {
	validationErrors := validate.NewErrors()

	err := validation.ValidateStruct(l,
		validation.Field(&l.Code, validation.Required),
		validation.Field(&l.Name, validation.Required),
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

func (l *Location) BeforeCreate(tx *pop.Connection) error {

	count, err := tx.Where("code = ?", l.Code).Count(l)

	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("code already in use")
	}

	validateName := &Location{
		Name: l.Name,
	}

	count, err = tx.Where("name = ?", l.Name).Count(validateName)

	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("name already in use")
	}

	return nil
}
