package services

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/models"
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

func (us *UserService) GetAll() (*models.Users, error) {
	users := &models.Users{}

	res := us.db.Find(users)

	return users, res.Error
}

func (us *UserService) Get(id uint) (*models.User, bool, error) {
	model := &models.User{}

	res := us.db.First(model, id)

	return model, model.ID != 0, res.Error
}

func (us *UserService) Update(model *models.User) error {
	res := us.db.Update(model)
	return res.Error
}

func (us *UserService) Delete(id uint) error {
	model := &models.Server{
		ID: id,
	}

	res := us.db.Delete(model)
	return res.Error
}

func (us *UserService) ChangePassword(id uint, newPass string) error {
	user, exists, err := us.Get(id)

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