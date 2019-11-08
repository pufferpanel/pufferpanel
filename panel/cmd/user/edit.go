package user

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/panel/database"
	"github.com/pufferpanel/pufferpanel/v2/panel/services"
	"github.com/spf13/cobra"
)

var EditUserCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a user",
	Run:   editUser,
	Args:  cobra.NoArgs,
}

func editUser(cmd *cobra.Command, args []string) {
	err := pufferpanel.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %s", err.Error())
		return
	}

	db, err := database.GetConnection()
	if err != nil {
		fmt.Printf("Error connecting to database: %s", err.Error())
		return
	}
	defer database.Close()

	var username string
	_ = survey.AskOne(&survey.Input{
		Message: "Username:",
	}, &username, survey.WithValidator(survey.Required))

	us := &services.User{DB: db}

	user, err := us.Get(username)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		fmt.Printf("No user with username '%s'", username)
		return
	} else if err != nil {
		fmt.Printf("Error getting user: %s", err.Error())
		return
	}

	action := ""
	_ = survey.AskOne(&survey.Select{
		Message: "Select option to edit",
		Options: []string{"Username", "Email", "Password", "Change Admin Status"},
	}, &action)

	switch action {
	case "Username":
		{
			prompt := ""
			_ = survey.AskOne(&survey.Input{
				Message: "New Username:",
			}, &prompt, survey.WithValidator(survey.Required))
			user.Username = prompt

			err = us.Update(user)
			if err != nil {
				fmt.Printf("Error updating username: %s", err.Error())
			}
		}
	case "Email":
		{
			prompt := ""
			_ = survey.AskOne(&survey.Input{
				Message: "New Email:",
			}, &prompt, survey.WithValidator(survey.Required))
			user.Email = prompt

			err = us.Update(user)
			if err != nil {
				fmt.Printf("Error updating email: %s", err.Error())
			}
		}
	case "Password":
		{
			prompt := ""
			_ = survey.AskOne(&survey.Password{
				Message: "New Password:",
			}, &prompt, survey.WithValidator(validatePassword))

			err = user.SetPassword(prompt)
			if err != nil {
				fmt.Printf("Error updating password: %s", err.Error())
			}

			err = us.Update(user)
			if err != nil {
				fmt.Printf("Error updating password: %s", err.Error())
			}
		}
	case "Change Admin Status":
		{
			prompt := false
			_ = survey.AskOne(&survey.Confirm{
				Message: "Set Admin Status: ",
			}, &prompt)

			ps := &services.Permission{DB: db}
			perms, err := ps.GetForUserAndServer(user.ID, nil)
			if err != nil {
				fmt.Printf("Error updating permissions: %s", err.Error())
				return
			}

			perms.Admin = prompt

			err = ps.UpdatePermissions(perms)
			if err != nil {
				fmt.Printf("Error updating password: %s", err.Error())
			}
		}
	}
}
