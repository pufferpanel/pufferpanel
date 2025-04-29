package main

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/pterm/pterm"
	"github.com/pufferpanel/pufferpanel/v3/database"
	"github.com/pufferpanel/pufferpanel/v3/groups"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/scopes"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"strings"
)

var AddUserCmd = &cobra.Command{
	Use:   "add",
	Short: "Add user",
	Run:   addUser,
	Args:  cobra.NoArgs,
}

var EditUserCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a user",
	Run:   editUser,
	Args:  cobra.NoArgs,
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
}

var addUsername string
var addEmail string
var addIsAdmin bool
var addPassword string

func init() {
	userCmd.AddCommand(AddUserCmd, EditUserCmd)

	AddUserCmd.Flags().StringVar(&addUsername, "name", "", "username")
	AddUserCmd.Flags().StringVar(&addEmail, "email", "", "email")
	AddUserCmd.Flags().BoolVar(&addIsAdmin, "admin", false, "if admin")
	AddUserCmd.Flags().StringVar(&addPassword, "password", "", "password")
}

func addUser(cmd *cobra.Command, args []string) {
	answers := userCreate{
		Username: addUsername,
		Email:    addEmail,
		Admin:    addIsAdmin,
		Password: addPassword,
	}

	//should we ask if this user is an admin should only appear if no flags are used
	promptAdmin := true
	if answers.Admin || answers.Username != "" || answers.Email != "" || answers.Password != "" {
		promptAdmin = false
	}

	questions := make([]*survey.Question, 0)

	if answers.Username == "" {
		questions = append(questions, &survey.Question{
			Name: "username",
			Prompt: &survey.Input{
				Message: "Username:",
			},
			Validate: validateUsername,
		})
	}

	if answers.Email == "" {
		questions = append(questions, &survey.Question{
			Name: "email",
			Prompt: &survey.Input{
				Message: "Email:",
			},
			Validate: validateEmail,
		})
	}

	if answers.Password == "" {
		questions = append(questions, &survey.Question{
			Name: "password",
			Prompt: &survey.Password{
				Message: "Password:",
			},
			Validate: validatePassword,
		})
	}

	if promptAdmin {
		questions = append(questions, &survey.Question{
			Name: "admin",
			Prompt: &survey.Confirm{
				Message: "Admin",
			},
		})
	}

	if len(questions) > 0 {
		_ = survey.Ask(questions, &answers)
	}

	db, err := database.GetConnection()
	if err != nil {
		fmt.Printf("Failed to connect to database: %s\n", err.Error())
		return
	}
	defer database.Close()

	if err := db.Transaction(func(tx *gorm.DB) error {
		user := &models.User{
			Username:       answers.Username,
			Email:          answers.Email,
			HashedPassword: "",
		}
		err = user.SetPassword(answers.Password)
		if err != nil {
			fmt.Printf("Failed to set password: %s\n", err.Error())
			return err
		}

		us := &services.User{DB: tx}
		err = us.Create(user)
		if err != nil {
			fmt.Printf("Failed to create user: %s\n", err.Error())
			return err
		}

		ps := &services.Permission{DB: tx}
		perms, err := ps.GetForUserAndServer(user.ID, "")
		if err != nil {
			fmt.Printf("Failed to get permissions: %s\n", err.Error())
			return err
		}

		if answers.Admin {
			perms.Scopes = scopes.AddScope(perms.Scopes, scopes.ScopeAdmin)
		}

		err = ps.UpdatePermissions(perms)
		if err != nil {
			fmt.Printf("Failed to apply permissions: %s\n", err.Error())
			return err
		}

		return nil
	}); err != nil {
		return
	}

	fmt.Printf("User added\n")
}

func validateEmail(val interface{}) error {
	email := val.(string)

	var viewModel models.UserView
	viewModel.Email = email
	err := viewModel.EmailValid(false)
	if err != nil {
		return err
	}

	return nil
}

func validateUsername(val interface{}) error {
	usr := val.(string)

	var viewModel models.UserView
	viewModel.Username = usr
	err := viewModel.UserNameValid(false)
	if err != nil {
		return err
	}

	return nil
}

func validatePassword(val interface{}) error {
	pw, ok := val.(string)
	if !ok {
		return errors.New("password is not a string")
	}

	us := &services.User{}
	if !us.IsSecurePassword(pw) {
		return errors.New("password must be at least 8 characters")
	}

	var secondAttempt string
	confirm := &survey.Password{
		Message: "Confirm Password",
	}
	_ = survey.AskOne(confirm, &secondAttempt)

	if secondAttempt != pw {
		return errors.New("password do not match")
	}

	return nil
}

type userCreate struct {
	Username string
	Email    string
	Password string
	Admin    bool
}

func editUser(cmd *cobra.Command, args []string) {
	if !groups.IsUserIn(groups.PufferPanelGroup) {
		fmt.Printf("You do not have permission to use this command")
		return
	}

	db, err := database.GetConnection()
	if err != nil {
		fmt.Printf("Error connecting to database: %s", err.Error())
		return
	}
	defer database.Close()

	var username string

	username, _ = pterm.DefaultInteractiveTextInput.WithDefaultText("Enter Username").Show()

	us := &services.User{DB: db}

	user, err := us.Get(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		pterm.Error.Printfln("No user with username '%s'\n", username)
		return
	} else if err != nil {
		pterm.Error.Printfln("Error getting user: %s\n", err.Error())
		return
	}

	var usernameAction = "Username"
	var emailAction = "Email"
	var passwordAction = "Password"
	var adminAction = "Admin Status"
	var remove2FAAction = "Remove 2FA"
	var quitAction = "Quit"

	var currentAction string

	for currentAction != quitAction {
		currentAction, _ = pterm.DefaultInteractiveSelect.WithOptions([]string{
			usernameAction,
			emailAction,
			passwordAction,
			adminAction,
			remove2FAAction,
			quitAction,
		}).WithFilter(false).WithMaxHeight(20).Show()

		switch currentAction {
		case usernameAction:
			{
				oldValue := user.Username
				user.Username, _ = pterm.DefaultInteractiveTextInput.WithDefaultText("Change username to").WithDefaultValue(oldValue).Show()

				err = us.Update(user)
				if err != nil {
					pterm.Error.Printfln("Error updating username: %s", err.Error())
				} else {
					pterm.Info.Printfln("Username updated from %s to %s", oldValue, user.Username)
				}
			}
		case emailAction:
			{
				oldValue := user.Username
				user.Email, _ = pterm.DefaultInteractiveTextInput.WithDefaultText("Change email to").WithDefaultValue(oldValue).Show()

				err = us.Update(user)
				if err != nil {
					pterm.Error.Printfln("Error updating email: %s", err.Error())
				} else {
					pterm.Info.Printfln("Email updated from %s to %s", oldValue, user.Username)
				}
			}
		case passwordAction:
			{
				password, _ := pterm.DefaultInteractiveTextInput.WithMask("*").WithDefaultText("Change password to").Show()

				err = user.SetPassword(password)
				if err != nil {
					pterm.Error.Printfln("Error updating password: %s", err.Error())
					break
				}
				err = us.Update(user)
				if err != nil {
					pterm.Error.Printfln("Error updating password: %s", err.Error())
				} else {
					pterm.Info.Printfln("Password updated")
				}
			}
		case adminAction:
			{
				result, _ := pterm.DefaultInteractiveContinue.WithDefaultText("Set as admin").WithOptions([]string{"yes", "no", "cancel"}).WithDefaultValue("cancel").Show()

				if result == "cancel" {
					break
				}

				ps := &services.Permission{DB: db}
				perms, err := ps.GetForUserAndServer(user.ID, "")
				if err != nil {
					pterm.Error.Printfln("Error updating permissions: %s", err.Error())
					break
				}

				//perms.Admin = prompt
				result = strings.ToLower(result)
				if result == "yes" || result == "y" {
					perms.Scopes = scopes.AddScope(perms.Scopes, scopes.ScopeAdmin)
				} else if result == "no" || result == "n" {
					perms.Scopes = scopes.RemoveScope(perms.Scopes, scopes.ScopeAdmin)
				} else {
					break
				}

				err = ps.UpdatePermissions(perms)
				if err != nil {
					pterm.Error.Printfln("Error updating password: %s", err.Error())
					break
				}

				if scopes.ContainsScope(perms.Scopes, scopes.ScopeAdmin) {
					pterm.Info.Printfln("Admin status added")
				} else {
					pterm.Info.Printfln("Admin status removed")
				}
			}
		case remove2FAAction:
			{
				result, _ := pterm.DefaultInteractiveConfirm.WithDefaultText("Remove 2FA").Show()
				if result {
					us := &services.User{DB: db}
					err = us.DisableOtp(user.ID)
					if err != nil {
						fmt.Printf("Error removing 2FA: %s", err.Error())
					} else {
						pterm.Info.Printfln("2FA removed")
					}
				}
			}
		}
	}
}
