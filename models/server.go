package models

import (
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/satori/go.uuid"
	"time"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Server struct {
	ID        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Uuid      uuid.UUID `json:"uuid" db:"uuid"`
	Name      string    `json:"name" db:"name"`
	User_ID   int       `json:"user_id" db:"user_id"`
	Node_ID   int       `json:"node_id" db:"node_id"`
}

// Validate gets run every time you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (s *Server) Validate(tx *pop.Connection) (*validate.Errors, error) {
	resultErrs := validate.NewErrors()

	err := validation.ValidateStruct(&s,
		validation.Field(&s.Uuid, validation.Required, is.UUID),
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.User_ID, validation.Required),
		validation.Field(&s.Node_ID, validation.Required),
	)

	errs := err.(validation.Errors)

	if err != nil && errs.Filter() != nil {
		for k, v := range errs {
			resultErrs.Add(k, v.Error())
		}
	}

	node := &Node{}
	err = tx.BelongsTo(s).All(&node)

	if err != nil {
		resultErrs.Add("node", err.Error())
	}

	if node == nil {
		resultErrs.Add("node", "node does not exist")
	}

	user := &User{}
	err = tx.BelongsTo(s).All(&user)

	if err != nil {
		resultErrs.Add("user", err.Error())
	}

	if user == nil {
		resultErrs.Add("user", "user does not exist")
	}

	return resultErrs, nil
}
