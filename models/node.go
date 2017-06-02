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
	Public_Ip   string  `json:"public_ip" db:"public_ip"`
	Private_Ip  string  `json:"private_ip" db:"private_ip"`
	Port        int     `json:"port" db:"port"`
}

// String is not required by pop and may be deleted
func (n Node) String() string {
	jn, _ := json.Marshal(n)
	return string(jn)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (n *Node) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}