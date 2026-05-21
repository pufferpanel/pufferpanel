package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"
	"github.com/pufferpanel/pufferpanel/v3/database"
	"github.com/pufferpanel/pufferpanel/v3/groups"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/scopes"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
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

	useFlags := false
	if answers.Username != "" || answers.Email != "" || answers.Password != "" {
		useFlags = true
	}

	firstAnswer := answers.Username == ""
	err := validateUsername(answers.Username)
	for err != nil {
		if !firstAnswer {
			pterm.Error.Println("Username validation failed: " + err.Error())
			if useFlags {
				os.Exit(1)
			}
		}
		firstAnswer = false
		answers.Username, _ = pterm.DefaultInteractiveTextInput.WithDefaultText("Username").Show()
		err = validateUsername(answers.Username)
	}

	firstAnswer = answers.Email == ""
	err = validateEmail(answers.Email)
	for err != nil {
		if !firstAnswer {
			pterm.Error.Println("Email validation failed: " + err.Error())
			if useFlags {
				os.Exit(1)
			}
		}
		firstAnswer = false
		answers.Email, _ = pterm.DefaultInteractiveTextInput.WithDefaultText("Email").Show()
		err = validateEmail(answers.Email)
	}

	firstAnswer = answers.Password == ""
	err = validatePassword(answers.Password)
	for err != nil {
		if !firstAnswer {
			pterm.Error.Println("Password validation failed: " + err.Error())
			if useFlags {
				os.Exit(1)
			}
		}
		firstAnswer = false
		answers.Password, _ = pterm.DefaultInteractiveTextInput.WithDefaultText("Password").WithMask("*").Show()
		err = validatePassword(answers.Password)
		if err != nil {
			continue
		}

		confirm, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Confirm Password").WithMask("*").Show()
		if answers.Password != confirm {
			err = errors.New("passwords do not match")
		}
	}

	if !useFlags {
		answers.Admin, _ = pterm.DefaultInteractiveConfirm.WithDefaultText("Set as Admin").Show()
	}

	db, err := database.GetConnection()
	if err != nil {
		pterm.Error.Printf("Failed to connect to database: %s\n", err.Error())
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
			pterm.Error.Printf("Failed to set password: %s\n", err.Error())
			return err
		}

		us := &services.User{DB: tx}
		err = us.Create(user)
		if err != nil {
			pterm.Error.Printf("Failed to create user: %s\n", err.Error())
			return err
		}

		ps := &services.Permission{DB: tx}
		perms, err := ps.GetForUserAndServer(user.ID, "")
		if err != nil {
			pterm.Error.Printf("Failed to get permissions: %s\n", err.Error())
			return err
		}

		if answers.Admin {
			perms.Scopes = scopes.AddScope(perms.Scopes, scopes.ScopeAdmin)
		}

		err = ps.UpdatePermissions(perms)
		if err != nil {
			pterm.Error.Printf("Failed to apply permissions: %s\n", err.Error())
			return err
		}

		return nil
	}); err != nil {
		return
	}

	pterm.Info.Printf("User added\n")
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
	return us.IsSecurePassword(pw)
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
				switch result {
				case "yes", "y":
					perms.Scopes = scopes.AddScope(perms.Scopes, scopes.ScopeAdmin)
				case "no", "n":
					perms.Scopes = scopes.RemoveScope(perms.Scopes, scopes.ScopeAdmin)
				default:
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
