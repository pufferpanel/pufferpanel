package models

import (
	"encoding/json"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/satori/go.uuid"
	"time"
)

type Node struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Name        string  `json:"name" db:"name"`
	Description string  `json:"description" db:"description"`
	Location_ID int     `json:"location_id" db:"location_id"`
}

// String is not required by pop and may be deleted
func (n Node) String() string {
	jn, _ := json.Marshal(n)
	return string(jn)
}

// Nodes is not required by pop and may be deleted
type Nodes []Node

// String is not required by pop and may be deleted
func (n Nodes) String() string {
	jn, _ := json.Marshal(n)
	return string(jn)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (n *Node) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (n *Node) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (n *Node) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
