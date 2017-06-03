package models

import (
	"encoding/json"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/satori/go.uuid"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/gobuffalo/buffalo/examples/html-crud/models"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type User struct {
	ID        int `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Uuid      uuid.UUID `json:"uuid" db:"uuid"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Language  string    `json:"language" db:"language"`
	Admin     bool      `json:"admin" db:"admin"`

	//private variable that's just backed by the database, we do not pass these outside this
	password  string    `json:"-" db:"password"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

func (u User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password)) == nil
}

func (u User) SetPassword(password string) error {
	//hash passwords using bcrypt for storage
	pw, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	u.password = string(pw)
	models.DB.ValidateAndSave(u)
	return err
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	resultErrs := validate.NewErrors()

	err := validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Uuid, validation.Required, is.UUID),
		validation.Field(&u.Username, validation.Required),
		validation.Field(&u.Language, validation.Required),
		validation.Field(&u.Admin, validation.Required),
		validation.Field(&u.password, validation.Required),
	)

	errs := err.(validation.Errors)

	if err != nil && errs.Filter() != nil {
		for k, v := range errs {
			resultErrs.Add(k, v.Error())
		}
	}

	return resultErrs, nil
}
