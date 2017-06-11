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
	UserId    int       `json:"user_id" db:"user_id"`
	NodeId    int       `json:"node_id" db:"node_id"`
}

type Servers []Server

func CreateServer() *Server {
	return &Server{
		Uuid: uuid.NewV4(),
	}
}

// Validate gets run every time you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (s *Server) Validate(tx *pop.Connection) (*validate.Errors, error) {
	resultErrs := validate.NewErrors()

	err := validation.ValidateStruct(&s,
		validation.Field(&s.Uuid, validation.Required, is.UUID),
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.UserId, validation.Required),
		validation.Field(&s.NodeId, validation.Required),
	)

	errs, ok := err.(validation.Errors)

	if ok && (err != nil && errs.Filter() != nil) {
		for k, v := range errs {
			resultErrs.Add(k, v.Error())
		}
	}

	node := &Node{}

	exists, err := tx.Where("id = ?", s.NodeId).Exists(node)

	if err != nil {
		resultErrs.Add("node", err.Error())
	}

	if !exists {
		resultErrs.Add("node", "node does not exist")
	}

	user := &User{}

	exists, err = tx.Where("id = ?", s.NodeId).Exists(user)

	if err != nil {
		resultErrs.Add("user", err.Error())
	}

	if !exists {
		resultErrs.Add("user", "user does not exist")
	}

	return resultErrs, nil
}
