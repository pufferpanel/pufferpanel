package grifts

import (
	"github.com/markbates/grift/grift"
	"github.com/markbates/pop"
	"github.com/pufferpanel/pufferpanel/models"
	"fmt"
	"errors"
)

var _ = grift.Add("db:seed", func(c *grift.Context) error {
	return models.DB.Transaction(func(tx *pop.Connection) error {
		location := &models.Location{
			Name: "seedlocation",
			Code: "seed1",
		}

		valErrs, err := tx.ValidateAndCreate(location)

		if valErrs != nil && len(valErrs.Errors) > 0 {
			fmt.Println(valErrs.Error())
			return errors.New(valErrs.Error())
		} else if err != nil {
			fmt.Println(err)
			return err
		}

		node := models.CreateNode()
		node.Name = "seednode"
		node.LocationId = location.ID
		valErrs, err = tx.ValidateAndCreate(node)

		if valErrs != nil && len(valErrs.Errors) > 0{
			fmt.Println(valErrs.Error())
			return errors.New(valErrs.Error())
		} else if err != nil {
			fmt.Println(err)
			return err
		}

		user := models.CreateUser()
		user.Username = "seeduser"
		user.Email = "seed@pufferpanel.com"
		user.SetPassword("seed")
		valErrs, err = tx.ValidateAndCreate(user)

		if valErrs != nil && len(valErrs.Errors) > 0 {
			fmt.Println(valErrs)
			return errors.New(valErrs.Error())
		} else if err != nil {
			fmt.Println(err)
			return err
		}

		server := models.CreateServer()
		server.Name = "seedserver"
		server.NodeId = node.ID
		server.UserId = user.ID
		valErrs, err = tx.ValidateAndCreate(server)

		if valErrs != nil && len(valErrs.Errors) > 0 {
			fmt.Println(valErrs)
			return errors.New(valErrs.Error())
		} else if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})
})
