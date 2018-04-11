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

type Client struct {
	ID           uuid.UUID `json:"id" db:"id"`
	HashedSecret string    `json:"-" db:"secret"`
	UserID       uuid.UUID `json:"user_id" db:"user_id" rw:"r"`
	Scopes       Scopes    `json:"scopes" has_many:"scopes"`
	Internal     bool      `json:"internal" db:"internal"`
	Description  string    `json:"description" db:"description"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Clients []Client

func GetClientsForUser(user *User, server *Server) (clients Clients, err error) {

	query := DB.Eager("scopes").BelongsTo(user)

	if server != nil {
		query.Where("scopes.server_id = ?", server.ID)
	}

	clients = Clients{}
	err = query.All(&clients)

	return
}

func (c *Client) GetScopesAsString() string {
	var result string
	for _, v := range c.Scopes {
		if result == "" {
			result = v.GetValue()
		} else {
			result += " " + v.GetValue()
		}
	}

	return result
}

func (c *Client) Validate(tx *pop.Connection) (*validate.Errors, error) {
	validationErrors := validate.NewErrors()

	err := validation.ValidateStruct(c,
		validation.Field(&c.Description, validation.Required),
		validation.Field(&c.HashedSecret, validation.Required),
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
