package models

import (
	"errors"
	"fmt"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"time"
)

type Node struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	Location   *Location  `json:"location" belongs_to:"location" fk_id:"location_id"`
	LocationID string 	  `json:"locationId" db:"location_id"`
	Code       string     `json:"code" db:"code"`
	Name       string     `json:"name" db:"name"`
	CreatedAt  time.Time  `json:"-" db:"created_at"`
	UpdatedAt  time.Time  `json:"-" db:"updated_at"`
}

type Nodes []Node

func CreateNode(location *Location, code, name string) (node Node, err error) {
	node = Node{
		Location: location,
		Name:     name,
		Code:     code,
	}

	return
}

func GetNodes() (nodes Nodes, err error) {
	nodes = Nodes{}
	err = DB.Eager().All(&nodes)
	return
}

func GetNodeByCode(code string) (node Node, err error) {
	node = Node{}
	query := DB.Where("code = ?", code)
	exists, err := query.Exists(&node)
	if exists {
		err = DB.Where("code = ?", code).First(&node)
	}
	return
}

func (n *Node) Delete() (err error) {
	err = DB.Destroy(n)
	return
}

func (n *Node) Save() (err error) {
	validationErrors, err := DB.ValidateAndSave(n)
	if validationErrors != nil && validationErrors.Count() > 0 {
		err = errors.New("model is invalid: " + validationErrors.Error())
	}
	return
}

func (n *Node) Validate(tx *pop.Connection) (*validate.Errors, error) {
	validationErrors := validate.NewErrors()

	err := validation.ValidateStruct(n,
		validation.Field(&n.Code, validation.Required),
		validation.Field(&n.Name, validation.Required),
	)
	errs, ok := err.(validation.Errors)

	if err == nil {
		ok = true
	}

	if n.Location == nil {
		validationErrors.Add("location", "location must be valid")
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

func (n *Node) BeforeCreate(tx *pop.Connection) error {

	count, err := tx.Where("code = ?", n.Code).Count(n)

	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("code already in use")
	}

	count, err = tx.Where("name = ?", n.Name).Count(n)

	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("name already in use")
	}

	return nil
}