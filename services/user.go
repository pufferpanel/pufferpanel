/*
 Copyright 2020 Padduck, LLC
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at
  	http://www.apache.org/licenses/LICENSE-2.0
  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package services

import (
	"bytes"
	"encoding/base64"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/config"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"image"
	"image/png"
	"strings"
)

type User struct {
	DB *gorm.DB
}

func (us *User) Get(username string) (*models.User, error) {
	model := &models.User{
		Username: username,
	}

	err := us.DB.Where(model).First(model).Error

	if err != nil {
		return nil, err
	}
	return model, nil
}

func (us *User) GetById(id uint) (*models.User, error) {
	model := &models.User{
		ID: id,
	}

	err := us.DB.Where(model).First(model).Error

	if err != nil {
		return nil, err
	}
	return model, nil
}

func (us *User) Login(email string, password string) (user *models.User, sessionToken string, otpNeeded bool, err error) {
	user = &models.User{
		Email: email,
	}

	err = us.DB.Where(user).First(user).Error

	if err != nil && gorm.ErrRecordNotFound != err {
		return
	}

	if user.ID == 0 || gorm.ErrRecordNotFound == err {
		err = pufferpanel.ErrInvalidCredentials
		return
	}

	if !us.IsValidCredentials(user, password) {
		err = pufferpanel.ErrInvalidCredentials
		return
	}

	if user.OtpActive {
		otpNeeded = true
		return
	}

	sessionToken, err = GenerateSession(user.ID)
	return
}

func (us *User) LoginOtp(email string, token string) (user *models.User, sessionToken string, err error) {
	user = &models.User{
		Email: email,
	}

	err = us.DB.Where(user).First(user).Error

	if err != nil && gorm.ErrRecordNotFound != err {
		return
	}

	if user.ID == 0 || gorm.ErrRecordNotFound == err {
		err = pufferpanel.ErrInvalidCredentials
		return
	}

	if !totp.Validate(token, user.OtpSecret) {
		err = pufferpanel.ErrInvalidCredentials
		return
	}

	sessionToken, err = GenerateSession(user.ID)
	return
}

func (us *User) IsValidCredentials(user *models.User, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)) == nil
}

func (us *User) GetByEmail(email string) (*models.User, error) {
	model := &models.User{
		Email: email,
	}

	err := us.DB.Where(model).First(model).Error

	if err != nil {
		return nil, err
	}
	return model, nil
}

func (us *User) Update(model *models.User) error {
	return us.DB.Save(model).Error
}

func (us *User) Delete(model *models.User) (err error) {
	return us.DB.Transaction(func(tx *gorm.DB) error {
		us.DB.Delete(models.Permissions{}, "user_id = ?", model.ID)
		us.DB.Delete(models.Client{}, "user_id = ?", model.ID)
		us.DB.Delete(models.User{}, "id = ?", model.ID)
		return nil
	})
}

func (us *User) Create(user *models.User) error {
	return us.DB.Create(user).Error
}

func (us *User) ChangePassword(username string, newPass string) error {
	user, err := us.Get(username)

	if err != nil {
		return err
	}

	err = user.SetPassword(newPass)
	if err != nil {
		return err
	}
	return us.Update(user)
}

func (us *User) GetOtpStatus(userId uint) (enabled bool, err error) {
	user, err := us.GetById(userId)
	if err != nil {
		return
	}

	enabled = user.OtpActive
	return
}

func (us *User) StartOtpEnroll(userId uint) (secret string, img string, err error) {
	user, err := us.GetById(userId)
	if err != nil {
		return
	}

	var key *otp.Key
	key, err = totp.Generate(totp.GenerateOpts{
		Issuer:      config.GetString("panel.settings.companyName"),
		AccountName: user.Email,
	})
	if err != nil {
		return
	}

	user.OtpSecret = key.Secret()
	user.OtpActive = false
	err = us.Update(user)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	var image image.Image
	image, err = key.Image(256, 256)
	if err != nil {
		return
	}
	png.Encode(&buf, image)
	img = "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())

	secret = key.Secret()
	return
}

func (us *User) ValidateOtpEnroll(userId uint, token string) error {
	user, err := us.GetById(userId)
	if err != nil {
		return err
	}

	if !totp.Validate(token, user.OtpSecret) {
		return pufferpanel.ErrInvalidCredentials
	}

	user.OtpActive = true
	return us.Update(user)
}

func (us *User) DisableOtp(userId uint, token string) error {
	user, err := us.GetById(userId)
	if err != nil {
		return err
	}

	if !totp.Validate(token, user.OtpSecret) {
		return pufferpanel.ErrInvalidCredentials
	}

	user.OtpSecret = ""
	user.OtpActive = false
	return us.Update(user)
}

func (us *User) Search(usernameFilter, emailFilter string, pageSize, page uint) (*models.Users, int64, error) {
	users := &models.Users{}

	query := us.DB

	usernameFilter = strings.Replace(usernameFilter, "*", "%", -1)
	emailFilter = strings.Replace(emailFilter, "*", "%", -1)

	if usernameFilter != "" && usernameFilter != "%" {
		query = query.Where("username LIKE ?", usernameFilter)
	}

	if emailFilter != "" && emailFilter != "%" {
		query = query.Where("email LIKE ?", emailFilter)
	}

	var count int64
	err := query.Model(users).Count(&count).Error

	if err != nil {
		return nil, 0, err
	}

	res := query.Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(users)

	return users, count, res.Error
}
