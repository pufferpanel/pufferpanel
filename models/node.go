package models

import (
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/satori/go.uuid"
	"time"
	"errors"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Node struct {
	ID          int `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Uuid        uuid.UUID `json:"uuid" db:"uuid"`
	Name        string  `json:"name" db:"name"`
	Description string  `json:"description" db:"description"`
	Location_ID int     `json:"location_id" db:"location_id"`
	Public_Ip   string  `json:"public_ip" db:"public_ip"`
	Private_Ip  string  `json:"private_ip" db:"private_ip"`
	Port        int     `json:"port" db:"port"`
}

type Nodes []Node

func CreateNode() *Node {
	return &Node {
		Uuid: uuid.NewV4(),
		Public_Ip: "127.0.0.1",
		Private_Ip: "127.0.0.1",
		Port: 5656,
	}
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (n *Node) Validate(tx *pop.Connection) (*validate.Errors, error) {
	resultErrs := validate.NewErrors()

	err := validation.ValidateStruct(&n,
		validation.Field(&n.Name, validation.Required),
		validation.Field(&n.Description),
		validation.Field(&n.Location_ID, validation.Required),
		validation.Field(&n.Public_Ip, validation.Required, is.Host),
		validation.Field(&n.Private_Ip, validation.Required, is.Host),
		validation.Field(&n.Port, validation.Required, is.Port),
	)

	errs, ok := err.(validation.Errors)

	if ok && (err != nil && errs.Filter() != nil) {
		for k, v := range errs {
			resultErrs.Add(k, v.Error())
		}
	}

	location := &Location{}
	err = tx.BelongsTo(n).All(&location)

	if err != nil {
		resultErrs.Add("location", err.Error())
	}

	if location == nil {
		resultErrs.Add("location", "location does not exist")
	}

	return resultErrs, nil
}

func (n *Node) BeforeDestroy(tx *pop.Connection) error {
	server := Server{}

	exists, err := tx.BelongsTo(n).Exists(&server)

	if err != nil {
		return err
	} else if exists {
		return errors.New("node is associated with servers")
	}

	return nil
}