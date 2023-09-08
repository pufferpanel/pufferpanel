package tests

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/database"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"gorm.io/gorm"
)

var loginNoLoginUser = &models.User{
	Username:  "loginNoLoginUser",
	Email:     "noscope@example.com",
	OtpActive: false,
}
var loginNoLoginUserPassword = "dontletmein"

var loginNoServerViewUser = &models.User{
	Username:  "loginNoServerViewUser",
	Email:     "test@example.com",
	OtpActive: false,
}
var loginNoServerViewUserPassword = "testing123"

var loginAdminUser = &models.User{
	Username:  "loginAdminUser",
	Email:     "admin@example.com",
	OtpActive: false,
}
var loginAdminUserPassword = "asdfasdf"

var loginNoAdminWithServersUser = &models.User{
	Username:  "loginNoAdminWithServersUser",
	Email:     "notadmin@example.com",
	OtpActive: false,
}
var loginNoAdminWithServersUserPassword = "dowiuzlaslf"

func init() {
	_ = loginNoLoginUser.SetPassword(loginNoLoginUserPassword)
	_ = loginNoServerViewUser.SetPassword(loginNoServerViewUserPassword)
	_ = loginAdminUser.SetPassword(loginAdminUserPassword)
	_ = loginNoAdminWithServersUser.SetPassword(loginNoAdminWithServersUserPassword)
}

func prepareUsers(db *gorm.DB) error {
	err := initNoLoginUser(db)
	if err != nil {
		return err
	}

	err = initLoginNoServersUser(db)
	if err != nil {
		return err
	}

	err = initLoginAdminUser(db)
	if err != nil {
		return err
	}

	err = initLoginNoAdminWithServersUser(db)
	if err != nil {
		return err
	}

	return nil
}

func initNoLoginUser(db *gorm.DB) error {
	return db.Create(loginNoLoginUser).Error
}

func initLoginNoServersUser(db *gorm.DB) error {
	err := db.Create(loginNoServerViewUser).Error
	if err != nil {
		return err
	}

	perms := &models.Permissions{
		UserId: &loginNoServerViewUser.ID,
		Scopes: []*pufferpanel.Scope{pufferpanel.ScopeLogin},
	}
	err = db.Create(perms).Error
	return err
}

func initLoginAdminUser(db *gorm.DB) error {
	err := db.Create(loginAdminUser).Error
	if err != nil {
		return err
	}

	perms := &models.Permissions{
		UserId: &loginAdminUser.ID,
		Scopes: []*pufferpanel.Scope{pufferpanel.ScopeAdmin},
	}
	err = db.Create(perms).Error
	return err
}

func initLoginNoAdminWithServersUser(db *gorm.DB) error {
	return db.Create(loginNoAdminWithServersUser).Error
}

func createSession(db *gorm.DB, user *models.User) (string, error) {
	ss := &services.Session{DB: db}
	return ss.CreateForUser(user)
}

func createSessionAdmin() (string, error) {
	db, err := database.GetConnection()
	if err != nil {
		return "", err
	}
	return createSession(db, loginAdminUser)
}
