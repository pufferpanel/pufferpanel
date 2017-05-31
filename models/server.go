package models

import (
	"encoding/json"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/satori/go.uuid"
	"time"
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

// String is not required by pop and may be deleted
func (s Server) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Servers is not required by pop and may be deleted
type Servers []Server

// String is not required by pop and may be deleted
func (s Servers) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (s *Server) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (s *Server) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (s *Server) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
