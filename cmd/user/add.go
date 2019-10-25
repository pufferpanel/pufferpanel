package user

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"github.com/spf13/cobra"
)

var AddUserCmd = &cobra.Command{
	Use:   "add",
	Short: "Add user",
	Run:   runAdd,
	Args:  cobra.NoArgs,
}

var addUsername string
var addEmail string
var addIsAdmin bool
var addPassword string

func init() {
	AddUserCmd.Flags().StringVar(&addUsername, "name", "", "username")
	AddUserCmd.Flags().StringVar(&addEmail, "email", "", "email")
	AddUserCmd.Flags().BoolVar(&addIsAdmin, "admin", false, "if admin")
	AddUserCmd.Flags().StringVar(&addPassword, "password", "", "password")
}

func runAdd(cmd *cobra.Command, args []string) {
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
			Validate: survey.Required,
		})
	}

	if answers.Email == "" {
		questions = append(questions, &survey.Question{
			Name: "email",
			Prompt: &survey.Input{
				Message: "Email:",
			},
			Validate: survey.Required,
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

	err := pufferpanel.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %s", err.Error())
		return
	}

	db, err := database.GetConnection()
	if err != nil {
		fmt.Printf("Failed to connect to database: %s", err.Error())
		return
	}
	defer database.Close()

	db = db.Begin()
	defer db.RollbackUnlessCommitted()

	user := &models.User{
		Username:       answers.Username,
		Email:          answers.Email,
		HashedPassword: "",
	}
	err = user.SetPassword(answers.Password)
	if err != nil {
		fmt.Printf("Failed to set password: %s", err.Error())
		return
	}

	us := &services.User{DB: db}
	err = us.Create(user)
	if err != nil {
		fmt.Printf("Failed to create user: %s", err.Error())
		return
	}

	ps := &services.Permission{DB: db}
	perms, err := ps.GetForUserAndServer(user.ID, nil)
	if err != nil {
		fmt.Printf("Failed to get permissions: %s", err.Error())
		return
	}
	perms.Admin = answers.Admin
	perms.ViewServer = true
	err = ps.UpdatePermissions(perms)
	if err != nil {
		fmt.Printf("Failed to apply permissions: %s", err.Error())
		return
	}

	err = db.Commit().Error
	if err != nil {
		fmt.Printf("Failed to save changes: %s", err.Error())
		return
	}

	fmt.Printf("User added")
}

func validatePassword(val interface{}) error {
	pw, ok := val.(string)
	if !ok || len(pw) < 6 {
		return errors.New("Password must be at least 6 characters")
	}
	var secondAttempt string
	confirm := &survey.Password{
		Message: "Confirm Password",
	}
	_ = survey.AskOne(confirm, &secondAttempt)

	if secondAttempt != pw {
		return errors.New("Passwords do not match")
	}

	return nil
}

type userCreate struct {
	Username string
	Email    string
	Password string
	Admin    bool
}
