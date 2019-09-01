package services

import (
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel"
	"github.com/pufferpanel/pufferpanel/models"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
)

type User struct {
	DB *gorm.DB
}

func (us *User) Get(username string) (*models.User, bool, error) {
	model := &models.User{
		Username: username,
	}

	res := us.DB.Where(model).FirstOrInit(model)

	return model, model.ID != 0, res.Error
}

func (us *User) Login(email string, password string) (sessionToken string, err error) {
	oauth2 := GetOAuth(us.DB)

	model := &models.User{
		Email: email,
	}

	err = us.DB.Where(model).First(model).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return
	}

	if model.ID == 0 || gorm.IsRecordNotFoundError(err) {
		err = pufferpanel.ErrInvalidCredentials
		return
	}

	if !us.IsValidCredentials(model, password) {
		err = pufferpanel.ErrInvalidCredentials
		return
	}

	sessionToken, err = oauth2.CreateSession(model)
	return
}

func (us *User) IsValidCredentials(user *models.User, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)) != nil
}

func (us *User) GetByEmail(email string) (*models.User, bool, error) {
	model := &models.User{
		Email: email,
	}

	res := us.DB.Where(model).FirstOrInit(model)

	return model, model.ID != 0, res.Error
}

func (us *User) Update(model *models.User) error {
	res := us.DB.Save(model)
	return res.Error
}

func (us *User) Delete(username string) error {
	model := &models.User{
		Username: username,
	}

	res := us.DB.Delete(model)
	return res.Error
}

func (us *User) Create(user *models.User) error {
	oauth2 := GetOAuth(us.DB)

	res := us.DB.Create(user)
	if res.Error != nil {
		return res.Error
	}

	name := ".internal_" + strconv.Itoa(int(user.ID))

	_, err := oauth2.Create(user, nil, name, true, "login")

	if err != nil {
		us.DB.Delete(user)
	}

	return err
}

func (us *User) ChangePassword(username string, newPass string) error {
	user, exists, err := us.Get(username)

	if err != nil {
		return err
	}

	if !exists {
		return pufferpanel.ErrUserNotFound
	}

	err = user.SetPassword(newPass)
	if err != nil {
		return err
	}
	return us.Update(user)
}

func (us *User) Search(usernameFilter, emailFilter string, pageSize, page uint) (*models.Users, uint, error) {
	users := &models.Users{}

	query := us.DB.Offset((page - 1) * pageSize).Limit(pageSize)

	usernameFilter = strings.Replace(usernameFilter, "*", "%", -1)
	emailFilter = strings.Replace(emailFilter, "*", "%", -1)

	if usernameFilter != "" && usernameFilter != "%" {
		query = query.Where("username LIKE ?", usernameFilter)
	}

	if emailFilter != "" && emailFilter != "%" {
		query = query.Where("email LIKE ?", emailFilter)
	}

	res := query.Find(users)

	var count uint
	err := query.Model(users).Count(&count).Error

	if err != nil {
		return nil, 0, err
	}

	return users, count, res.Error
}
