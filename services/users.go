package services

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/models"
	"strings"
)

type UserService struct {
	db *gorm.DB
}

func GetUserService() (*UserService, error) {
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	service := &UserService{
		db: db,
	}

	return service, nil
}

func (us *UserService) Get(username string) (*models.User, bool, error) {
	model := &models.User{
		Username: username,
	}

	res := us.db.Where(model).FirstOrInit(model)

	return model, model.ID != 0, res.Error
}

func (us *UserService) GetByEmail(email string) (*models.User, bool, error) {
	model := &models.User{
		Email: email,
	}

	res := us.db.Where(model).FirstOrInit(model)

	return model, model.ID != 0, res.Error
}

func (us *UserService) Update(model *models.User) error {
	res := us.db.Save(model)
	return res.Error
}

func (us *UserService) Delete(username string) error {
	model := &models.User{
		Username: username,
	}

	res := us.db.Delete(model)
	return res.Error
}

func (us *UserService) Create(user *models.User) error {
	res := us.db.Create(user)
	return res.Error
}

func (us *UserService) ChangePassword(username string, newPass string) error {
	user, exists, err := us.Get(username)

	if err != nil {
		return err
	}

	if !exists {
		return errors.New("no such user")
	}

	err = user.SetPassword(newPass)
	if err != nil {
		return err
	}
	return us.Update(user)
}

func (us *UserService) Search (usernameFilter, emailFilter string, pageSize, page uint) (*models.Users, error) {
	users := &models.Users{}

	query := us.db.Offset((page - 1) * pageSize).Limit(pageSize)

	usernameFilter = strings.Replace(usernameFilter, "*", "%", -1)
	emailFilter = strings.Replace(emailFilter, "*", "%", -1)

	if usernameFilter != "" && usernameFilter != "%" {
		query = query.Where("username LIKE ?", usernameFilter)
	}

	if emailFilter != "" && emailFilter != "%" {
		query = query.Where("email LIKE ?", emailFilter)
	}

	res := query.Find(users)

	return users, res.Error
}