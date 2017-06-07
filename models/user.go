package models

import (
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/satori/go.uuid"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/gobuffalo/buffalo/examples/html-crud/models"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"errors"
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

type Users []User

func CreateUser() *User {
	return &User{
		Uuid: uuid.NewV4(),
		Language: "en_us",
		Admin: false,
	}
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

	errs, ok := err.(validation.Errors)

	if ok && (err != nil && errs.Filter() != nil) {
		for k, v := range errs {
			resultErrs.Add(k, v.Error())
		}
	}

	return resultErrs, nil
}

func (u *User) BeforeDestroy(tx *pop.Connection) error {
	server := Server{}

	exists, err := tx.BelongsTo(u).Exists(&server)

	if err != nil {
		return err
	} else if exists {
		return errors.New("user is associated with servers")
	}

	return nil
}