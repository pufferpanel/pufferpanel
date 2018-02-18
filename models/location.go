package models

import (
	"errors"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/satori/go.uuid"
	"github.com/go-ozzo/ozzo-validation"
	"fmt"
)

type Location struct {
	ID        uuid.UUID `json:"id" db:id`
	Code      string    `json:"code" db:"code"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

type Locations []Location

func GetLocations() (locations Locations, err error) {
	locations = Locations{}
	err = DB.All(locations)
	return
}

func GetLocationById(id string) (location Location, err error) {
	location = Location{}
	err = DB.Find(&location, id)
	return
}

func GetLocationByCode(code string) (location Location, err error) {
	location = Location{}
	err = DB.Where("code = ?", code).First(&location)
	return
}

func CreateLocation(code, name string) (location Location, err error) {
	id, err := uuid.NewV4()
	location = Location{
		ID: id,
		Code: code,
		Name: name,
	}

	err = DB.Create(location)
	return
}

func (l *Location) Delete() (err error) {
	err = DB.Destroy(l)
	return
}

func (l *Location) Save() (err error) {
	err = DB.Save(l)
	return
}

func (l *Location) Validate(tx *pop.Connection) (*validate.Errors, error) {
	validationErrors := validate.NewErrors()

	if l.ID == uuid.Nil {
		validationErrors.Add("id", "id is required")
	}

	err := validation.ValidateStruct(l,
		validation.Field(&l.ID, validation.Required),
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
	validateCode := &Location{
		Code: l.Code,
	}

	count, err := tx.Count(validateCode)

	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("code already in use")
	}

	validateName := &Location{
		Name: l.Name,
	}

	count, err = tx.Count(validateName)

	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("name already in use")
	}

	return nil
}
