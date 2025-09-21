package services

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"image"
	"image/png"
	"math/big"
	"net/http"
	"strings"

	webauthnProto "github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
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

type PasswordLoginResult struct {
	User *models.User `json:"-"`
	NeedsSecondFactor bool `json:"needsSecondFactor"`
	OtpEnabled bool `json:"otpEnabled"`
	WebauthnChallenge *webauthnProto.CredentialAssertion `json:"webauthnChallenge"`
	WebauthnSessionData *webauthn.SessionData `json:"-"`
}

func (us *User) ValidatePasswordLogin(email string, password string) (result PasswordLoginResult, err error) {
	result.User = &models.User{
		Email: email,
	}

	err = us.DB.Where(result.User).First(result.User).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	if result.User.ID == 0 || errors.Is(err, gorm.ErrRecordNotFound) {
		err = pufferpanel.ErrInvalidCredentials
		return
	}

	if !us.IsValidCredentials(result.User, password) {
		err = pufferpanel.ErrInvalidCredentials
		return
	}

	if result.User.OtpActive {
		result.NeedsSecondFactor = true
		result.OtpEnabled = true
	}

	passkeys, err := us.getPasskeys(result.User.ID)
	if err != nil {
		return
	}

	if len(passkeys) > 0 {
		result.NeedsSecondFactor = true
		w, err := us.webAuthn()
		if err != nil {
			return result, err
		}

		assertion, sessionData, err := w.BeginMediatedLogin(
			&models.WebAuthnUser{
				User: result.User,
				Credentials: passkeys,
			},
			webauthnProto.MediationDefault,
		)
		result.WebauthnChallenge = assertion
		result.WebauthnSessionData = sessionData
		return result, err
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
		return us.Update(user)
	})
}

func (us *User) webAuthn() (*webauthn.WebAuthn, error) {
	rpConfig := &webauthn.Config{
		RPDisplayName: config.CompanyName.Value(),
		RPID: strings.TrimPrefix(config.MasterUrl.Value(), "https://"),
		RPOrigins: []string{config.MasterUrl.Value()},
	}

	return webauthn.New(rpConfig)
}

func (us *User) getPasskeys(userId uint) ([]webauthn.Credential, error) {
	var credentials []*models.WebauthnCredential
	model := &models.WebauthnCredential{UserId: userId}
	err := us.DB.Where(model).Find(&credentials).Error
	if err != nil {
		return nil, err
	}

	result := make([]webauthn.Credential, len(credentials))
	for k, v := range credentials {
		result[k] = v.Credential
	}

	return result, nil
}

func (us *User) GetPasskeyDescriptors(userId uint) ([]models.WebauthnCredentialView, error) {
	var credentials []*models.WebauthnCredential
	model := &models.WebauthnCredential{UserId: userId}
	err := us.DB.Where(model).Find(&credentials).Error
	if err != nil {
		return nil, err
	}

	views := make([]models.WebauthnCredentialView, len(credentials))

	for k, v := range credentials {
		views[k] = models.WebauthnCredentialView{
			ID: v.ID,
			Name: v.Name,
			Descriptor: v.Credential.Descriptor(),
		}
	}

	return views, nil
}

func (us *User) StartPasskeyEnroll(userId uint) (*webauthnProto.CredentialCreation, *webauthn.SessionData, error) {
	w, err := us.webAuthn()
	if err != nil {
		return nil, nil, err
	}

	user, err := us.GetById(userId)
	if err != nil {
		return nil, nil, err
	}

	knownCredentials, err := us.getPasskeys(user.ID)
	if err != nil {
		return nil, nil, err
	}

	c, s, err := w.BeginMediatedRegistration(
		&models.WebAuthnUser{
			User: user,
			Credentials: knownCredentials,
		},
		webauthnProto.MediationDefault,
		webauthn.WithExclusions(webauthn.Credentials(knownCredentials).CredentialDescriptors()),
		webauthn.WithExtensions(map[string]any{"credProps": true}),
	)
	return c, s, err
}

func (us *User) ValidatePasskeyEnroll(userId uint, name string, request *http.Request, sessionData webauthn.SessionData) error {
	w, err := us.webAuthn()
	if err != nil {
		return err
	}

	user, err := us.GetById(userId)
	if err != nil {
		return err
	}

	knownCredentials, err := us.getPasskeys(userId)
	if err != nil {
		return err
	}

	credential, err := w.FinishRegistration(
		&models.WebAuthnUser{
			User: user,
			Credentials: knownCredentials,
		},
		sessionData,
		request,
	)
	if err != nil {
		return err
	}

	model := &models.WebauthnCredential{
		ID: credential.ID,
		Name: name,
		UserId: user.ID,
		Credential: *credential,
	}

	return us.DB.Create(model).Error
}

func (us *User) StartPasskeyLogin(email string) (*webauthnProto.CredentialAssertion, *webauthn.SessionData, error) {
	w, err := us.webAuthn()
	if err != nil {
		return nil, nil, err
	}

	user, err := us.GetByEmail(email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, pufferpanel.ErrInvalidCredentials
	} else if err != nil {
		return nil, nil, err
	}

	if !user.AllowPasswordlessLogin {
		return nil, nil, pufferpanel.ErrInvalidCredentials
	}

	knownCredentials, err := us.getPasskeys(user.ID)
	if err != nil {
		return nil, nil, err
	}

	if len(knownCredentials) == 0 {
		return nil, nil, pufferpanel.ErrInvalidCredentials
	}

	return w.BeginMediatedLogin(
		&models.WebAuthnUser{
			User: user,
			Credentials: knownCredentials,
		},
		webauthnProto.MediationDefault,
	)
}

func (us *User) ValidatePasskeyLogin(email string, request *http.Request, sessionData webauthn.SessionData) (*models.User, error) {
	w, err := us.webAuthn()
	if err != nil {
		return nil, err
	}

	user, err := us.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	credentials, err := us.getPasskeys(user.ID)
	if err != nil {
		return nil, err
	}

	credential, err := w.FinishLogin(
		&models.WebAuthnUser{
			User: user,
			Credentials: credentials,
		},
		*&sessionData,
		request,
	)
	if err != nil {
		return nil, err
	}

	model := &models.WebauthnCredential{
		ID: credential.ID,
		UserId: user.ID,
	}

	return user, us.DB.Model(model).Updates(models.WebauthnCredential{Credential: *credential}).Error
}

func (us *User) DeletePasskey(id string) error {
	data := []byte(id)
	out := make([]byte, base64.RawURLEncoding.DecodedLen(len(data)))
	_, err := base64.RawURLEncoding.Decode(out, data)
	if err != nil {
		return err
	}

	model := &models.WebauthnCredential{
		ID: out,
	}

	return us.DB.Delete(model).Error
}

func (us *User) SetPasswordlessLogin(userId uint, value bool) error {
	user, err := us.GetById(userId)
	if err != nil {
		return err
	}

	user.AllowPasswordlessLogin = value
	return us.DB.Save(user).Error
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
