package services

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"image"
	"image/png"
	"math/big"
	"strings"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

func (us *User) ValidateLogin(email string, password string) (user *models.User, otpNeeded bool, err error) {
	user = &models.User{
		Email: email,
	}

	err = us.DB.Where(user).First(user).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	if user.ID == 0 || errors.Is(err, gorm.ErrRecordNotFound) {
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
	return
}

func (us *User) ValidOtp(email string, token string) (user *models.User, err error) {
	user = &models.User{
		Email: email,
	}

	err = us.DB.Where(user).First(user).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	if user.ID == 0 || errors.Is(err, gorm.ErrRecordNotFound) {
		err = pufferpanel.ErrInvalidCredentials
		return
	}

	if !totp.Validate(token, user.OtpSecret) {
		// check recovery codes
		rc := &models.RecoveryCode{
			UserId: user.ID,
		}
		rc.SetCode(token)

		res := us.DB.Where(rc).Delete(rc)
		if res.Error != nil {
			err = res.Error
		}
		if res.RowsAffected == 0 {
			err = pufferpanel.ErrInvalidCredentials
		}
	}
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
		tx.Delete(models.Permissions{}, "user_id = ?", model.ID)
		tx.Delete(models.Client{}, "user_id = ?", model.ID)
		tx.Delete(models.Session{}, "user_id = ?", model.ID)
		tx.Delete(models.User{}, "id = ?", model.ID)
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

func (us *User) StartOtpEnroll(userId uint) (secret string, imgStr string, err error) {
	user, err := us.GetById(userId)
	if err != nil {
		return
	}

	var key *otp.Key
	key, err = totp.Generate(totp.GenerateOpts{
		Issuer:      config.CompanyName.Value(),
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
	var img image.Image
	img, err = key.Image(256, 256)
	if err != nil {
		return
	}
	png.Encode(&buf, img)

	imgStr = "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	secret = key.Secret()
	return
}

func (us *User) gerateOtpRecoveryCode(length int) (string, error) {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	code := make([]byte, length)
	for i := range length {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		code[i] = chars[num.Int64()]
	}
	return string(code), nil
}

func (us *User) generateOtpRecoveryCodes(numCodes int) ([]string, error) {
	codes := make([]string, numCodes)
	for i := range numCodes {
		code, err := us.gerateOtpRecoveryCode(12)
		if err != nil {
			return nil, err
		}
		codes[i] = code
	}
	return codes, nil
}

func (us *User) ValidateOtpEnroll(userId uint, token string) ([]string, error) {
	user, err := us.GetById(userId)
	if err != nil {
		return nil, err
	}

	if !totp.Validate(token, user.OtpSecret) {
		return nil, pufferpanel.ErrInvalidCredentials
	}

	codes, err := us.generateOtpRecoveryCodes(10)
	if err != nil {
		return nil, err
	}

	rcs := make([]models.RecoveryCode, len(codes))
	for i, code := range codes {
		rcs[i] = models.RecoveryCode{
			UserId: user.ID,
		}
		rcs[i].SetCode(code)
	}

	err = us.DB.Create(&rcs).Error
	if err != nil {
		return nil, err
	}

	user.OtpActive = true
	return codes, us.Update(user)
	
}

func (us *User) RegenerateOtpRecoveryCodes(userId uint) ([]string, error) {
	user, err := us.GetById(userId)
	if err != nil {
		return nil, err
	}

	codes, err := us.generateOtpRecoveryCodes(10)
	if err != nil {
		return nil, err
	}

	rcs := make([]models.RecoveryCode, len(codes))
	for i, code := range codes {
		rcs[i] = models.RecoveryCode{
			UserId: user.ID,
		}
		rcs[i].SetCode(code)
	}

	err = us.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&models.RecoveryCode{}, "user_id = ?", user.ID).Error
		if err != nil {
			return err
		}

		return tx.Create(&rcs).Error
	})

	return codes, err
}

func (us *User) DisableOtp(userId uint) error {
	user, err := us.GetById(userId)
	if err != nil {
		return err
	}

	return us.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&models.RecoveryCode{}, "user_id = ?", user.ID).Error
		if err != nil {
			return err
		}

		user.OtpSecret = ""
		user.OtpActive = false
		return tx.Save(user).Error
	})
}

func (us *User) Search(usernameFilter, emailFilter string, pageSize, page uint) ([]*models.User, int64, error) {
	var users []*models.User

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

	res := query.Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&users)

	return users, count, res.Error
}

func (us *User) IsSecurePassword(password string) error {
	//TODO: Change to use validator
	if len(password) < 8 {
		return pufferpanel.ErrFieldLength("password", 8, 72)
	}
	return nil
}
