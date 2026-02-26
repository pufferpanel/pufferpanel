package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/pufferpanel/pufferpanel/v3/scopes"
	"github.com/pufferpanel/pufferpanel/v3/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	uuid "github.com/gofrs/uuid/v5"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
)

func registerSelf(g *gin.RouterGroup) {
	g.Handle("GET", "", middleware.RequiresPermission(scopes.ScopeLogin), getSelf)
	g.Handle("PUT", "", middleware.RequiresPermission(scopes.ScopeSelfEdit), updateSelf)
	g.Handle("OPTIONS", "", response.CreateOptions("GET", "PUT"))

	g.Handle("GET", "/otp", middleware.RequiresPermission(scopes.ScopeSelfEdit), getOtpStatus)
	g.Handle("POST", "/otp", middleware.RequiresPermission(scopes.ScopeSelfEdit), startOtpEnroll)
	g.Handle("PUT", "/otp", middleware.RequiresPermission(scopes.ScopeSelfEdit), validateOtpEnroll)
	g.Handle("OPTIONS", "/otp", response.CreateOptions("GET", "POST", "PUT"))

	g.Handle("POST", "/otp/recovery", middleware.RequiresPermission(scopes.ScopeSelfEdit), regenerateOtpRecoveryCodes)
	g.Handle("DELETE", "/otp/:token", middleware.RequiresPermission(scopes.ScopeSelfEdit), disableOtp)
	g.Handle("OPTIONS", "/otp/:token", response.CreateOptions("POST", "DELETE"))

	g.Handle("GET", "/passkey", middleware.RequiresPermission(scopes.ScopeSelfEdit), getPasskeys)
	g.Handle("POST", "/passkey", middleware.RequiresPermission(scopes.ScopeSelfEdit), startPasskeyEnroll)
	g.Handle("PUT", "/passkey", middleware.RequiresPermission(scopes.ScopeSelfEdit), validatePasskeyEnroll)
	g.Handle("OPTIONS", "/passkey", response.CreateOptions("GET", "OPTIONS", "PUT"))

	g.Handle("DELETE", "/passkey/:id", middleware.RequiresPermission(scopes.ScopeSelfEdit), deletePasskey)
	g.Handle("OPTIONS", "/passkey/:id", response.CreateOptions("DELETE"))

	g.Handle("PUT", "/passwordless/:value", middleware.RequiresPermission(scopes.ScopeSelfEdit), setPasswordlessLogin)
	g.Handle("OPTIONS", "/passwordless/:value", response.CreateOptions("PUT"))

	g.Handle("GET", "/oauth2", middleware.RequiresPermission(scopes.ScopeSelfClients), getPersonalOAuth2Clients)
	g.Handle("POST", "/oauth2", middleware.RequiresPermission(scopes.ScopeSelfClients), createPersonalOAuth2Client)
	g.Handle("OPTIONS", "/oauth2", response.CreateOptions("GET", "POST"))

	g.Handle("DELETE", "/oauth2/:clientId", middleware.RequiresPermission(scopes.ScopeSelfClients), deletePersonalOAuth2Client)
	g.Handle("OPTIONS", "/oauth2/:clientId", response.CreateOptions("DELETE"))

	g.Handle("GET", "/publickey", middleware.RequiresPermission(scopes.ScopeSelfEdit), getPublicKeys)
	g.Handle("PUT", "/publickey", middleware.RequiresPermission(scopes.ScopeSelfEdit), addPublicKey)
	g.Handle("DELETE", "/publickey", middleware.RequiresPermission(scopes.ScopeSelfEdit), deletePublicKey)
	g.Handle("OPTIONS", "/publickey", response.CreateOptions("GET", "PUT", "DELETE"))
}

// @Summary Get your user info
// @Description Gets the user information of the current user
// @Success 200 {object} models.UserView
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Router /api/self [get]
// @Security OAuth2Application[login]
func getSelf(c *gin.Context) {
	user := getUserFromContext(c)

	c.JSON(http.StatusOK, models.FromUser(user))
}

// @Summary Update your user
// @Description Update user information for your current user
// @Success 204 {object} nil
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Param user body models.UserView true "User information"
// @Router /api/self [PUT]
// @Security OAuth2Application[self.edit]
func updateSelf(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	user := getUserFromContext(c)

	var viewModel models.UserView
	if err := c.BindJSON(&viewModel); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	if err := viewModel.Valid(true); response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	if viewModel.Password == "" {
		response.HandleError(c, pufferpanel.ErrFieldRequired("password"), http.StatusBadRequest)
		return
	}

	if !us.IsValidCredentials(user, viewModel.Password) {
		response.HandleError(c, pufferpanel.ErrInvalidCredentials, http.StatusInternalServerError)
		return
	}

	var oldEmail string
	if user.Email != viewModel.Email {
		oldEmail = user.Email
	}

	viewModel.CopyToModel(user)

	passwordChanged := false
	if viewModel.NewPassword != "" {
		if us.IsSecurePassword(viewModel.NewPassword) != nil {
			response.HandleError(c, pufferpanel.ErrPasswordRequirements, http.StatusBadRequest)
			return
		}

		passwordChanged = true
		err := user.SetPassword(viewModel.NewPassword)
		if response.HandleError(c, err, http.StatusInternalServerError) {
			return
		}
	}

	if err := us.Update(user); response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	if oldEmail != "" {
		err := services.GetEmailService().SendEmail(oldEmail, "emailChanged", map[string]interface{}{
			"NEW_EMAIL": user.Email,
		}, true)
		if err != nil {
			logging.Error.Printf("Error sending email: %s\n", err)
		}
	}

	if passwordChanged {
		err := services.GetEmailService().SendEmail(user.Email, "passwordChanged", nil, true)
		if err != nil {
			logging.Error.Printf("Error sending email: %s\n", err)
		}
	}

	c.Status(http.StatusNoContent)
}

func getOtpStatus(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	user := getUserFromContext(c)

	otpEnabled, err := us.GetOtpStatus(user.ID)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"otpEnabled": otpEnabled,
	})
}

func startOtpEnroll(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	user := getUserFromContext(c)

	secret, img, err := us.StartOtpEnroll(user.ID)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"secret": secret,
		"img":    img,
	})
}

func validateOtpEnroll(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	user := getUserFromContext(c)

	request := &ValidateOtpRequest{}

	err := c.BindJSON(request)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	recoveryCodes, err := us.ValidateOtpEnroll(user.ID, request.Token)
	if errors.Is(err, pufferpanel.ErrInvalidCredentials) {
		response.HandleError(c, err, http.StatusBadRequest)
		return
	}
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = services.GetEmailService().SendEmail(user.Email, "otpEnabled", nil, true)
	if err != nil {
		logging.Error.Printf("Error sending email: %s\n", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"recoveryCodes": recoveryCodes,
	})
}

func regenerateOtpRecoveryCodes(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	user := getUserFromContext(c)

	request := &ValidateOtpRequest{}

	err := c.BindJSON(request)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	_, err = us.ValidOtp(user.Email, request.Token)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	recoveryCodes, err := us.RegenerateOtpRecoveryCodes(user.ID)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = services.GetEmailService().SendEmail(user.Email, "recoveryCodesRegenerated", nil, true)
	if err != nil {
		logging.Error.Printf("Error sending email: %s\n", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"recoveryCodes": recoveryCodes,
	})
}

func disableOtp(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	user := getUserFromContext(c)
	token := c.Param("token")

	_, err := us.ValidOtp(user.Email, token)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	err = us.DisableOtp(user.ID)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = services.GetEmailService().SendEmail(user.Email, "otpDisabled", nil, true)
	if err != nil {
		logging.Error.Printf("Error sending email: %s\n", err)
	}
	c.Status(http.StatusNoContent)
}

// @Summary Get registered passkeys
// @Description Gets descriptors of the known passkeys the logged-in user has registered
// @Success 200 {object} []models.WebauthnCredentialView
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Router /api/self/passkey [GET]
// @Security OAuth2Application[self.edit]
func getPasskeys(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	user := getUserFromContext(c)

	credentials, err := us.GetPasskeyDescriptors(user.ID)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, credentials)
}

// @Summary Start registering a new passkey
// @Success 200 {object} protocol.CredentialCreation
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Param request body EnrollPasskeyRequest true "Information for the passkey to register"
// @Router /api/self/passkey [POST]
// @Security OAuth2Application[self.edit]
func startPasskeyEnroll(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	user := getUserFromContext(c)

	var request EnrollPasskeyRequest
	err := c.BindJSON(&request)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	creation, sessionData, err := us.StartPasskeyEnroll(user.ID)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	sessionDataJson, err := json.Marshal(sessionData)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	userSession := sessions.Default(c)
	userSession.Set("passkeyName", request.Name)
	userSession.Set("passkeyEnroll", sessionDataJson)
	err = userSession.Save()
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, creation)
}

// @Summary Complete registering a new passkey
// @Success 204 {object} nil
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Param request body protocol.CredentialCreationResponse true "The response object for the passkeys registration challenge"
// @Router /api/self/passkey [PUT]
// @Security OAuth2Application[self.edit]
func validatePasskeyEnroll(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	user := getUserFromContext(c)

	sessionData := webauthn.SessionData{}
	userSession := sessions.Default(c)
	name := userSession.Get("passkeyName").(string)
	sessionDataJson := userSession.Get("passkeyEnroll").([]byte)
	err := json.Unmarshal(sessionDataJson, &sessionData)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = us.ValidatePasskeyEnroll(user.ID, name, c.Request, sessionData)
	if errors.Is(err, pufferpanel.ErrInvalidCredentials) {
		response.HandleError(c, err, http.StatusBadRequest)
		return
	}
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Remove a passkey
// @Success 204 {object} nil
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Router /api/self/passkey/:id [DELETE]
// @Security OAuth2Application[self.edit]
func deletePasskey(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	id := c.Param("id")

	err := us.DeletePasskey(id)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Set if passwordless login is allowed for the current user or not
// @Success 204 {object} nil
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Router /api/self/passwordless/:value [PUT]
// @Security OAuth2Application[self.edit]
func setPasswordlessLogin(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	user := getUserFromContext(c)
	value, err := strconv.ParseBool(c.Param("value"))
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	err = us.SetPasswordlessLogin(user.ID, value)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Gets registered OAuth2 clients
// @Description Gets known OAuth2 clients the logged-in user has registered
// @Success 200 {object} []models.Client
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Router /api/self/oauth2 [GET]
// @Security OAuth2Application[self.clients]
func getPersonalOAuth2Clients(c *gin.Context) {
	user := getUserFromContext(c)

	db := middleware.GetDatabase(c)
	os := &services.OAuth2{DB: db}

	clients, err := os.GetForUser(user.ID)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, &clients)
}

// @Summary Create an account-level OAuth2 client
// @Success 200 {object} models.Client
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Param client body models.Client false "Information for the client to create"
// @Router /api/self/oauth2 [POST]
// @Security OAuth2Application[self.clients]
func createPersonalOAuth2Client(c *gin.Context) {
	user := getUserFromContext(c)

	db := middleware.GetDatabase(c)
	os := &services.OAuth2{DB: db}

	var request models.Client
	err := c.BindJSON(&request)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	id, err := uuid.NewV4()
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}
	client := &models.Client{
		ClientId:    id.String(),
		UserId:      user.ID,
		Name:        request.Name,
		Description: request.Description,
	}

	secret, err := utils.GenerateRandomString(36)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	client.ClientSecret = secret

	err = client.SetClientSecret(secret)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = os.Create(client)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = services.GetEmailService().SendEmail(user.Email, "oauthCreated", nil, true)
	if err != nil {
		logging.Error.Printf("Error sending email: %s\n", err)
	}

	c.JSON(http.StatusOK, client)
}

// @Summary Deletes an account-level OAuth2 client
// @Success 204 {object} nil
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Param id path string true "Information for the client to create"
// @Router /api/self/oauth2/{id} [DELETE]
// @Security OAuth2Application[self.clients]
func deletePersonalOAuth2Client(c *gin.Context) {
	user := getUserFromContext(c)
	clientId := c.Param("clientId")

	db := middleware.GetDatabase(c)
	os := &services.OAuth2{DB: db}

	client, err := os.Get(clientId)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	//ensure the client id is specific for this server, and this user
	if client.UserId != user.ID {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	err = os.Delete(client.ClientId)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = services.GetEmailService().SendEmail(user.Email, "oauthDeleted", nil, true)
	if err != nil {
		logging.Error.Printf("Error sending email: %s\n", err)
	}
	c.Status(http.StatusNoContent)
}

// @Summary Gets public keys for the current user
// @Success 200 {object} models.PublicKeyView
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Router /api/self/publickey [GET]
// @Security OAuth2Application[self.edit]
func getPublicKeys(c *gin.Context) {
	user := getUserFromContext(c)
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	keys, err := us.GetPublicKeysForUser(user.ID)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, keys.ToView())
}

// @Summary Adds the given key to the current user
// @Success 204 {object} nil
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Router /api/self/publickey [PUT]
// @Security OAuth2Application[self.edit]
func addPublicKey(c *gin.Context) {
	user := getUserFromContext(c)
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	keys := &KeyRequest{}
	err := c.BindJSON(&keys)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = us.AddPublicKey(user.ID, keys.Key)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Removes the given key to the current user
// @Success 204 {object} nil
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 403 {object} pufferpanel.ErrorResponse
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Router /api/self/publickey [DELETE]
// @Security OAuth2Application[self.edit]
func deletePublicKey(c *gin.Context) {
	user := getUserFromContext(c)
	db := middleware.GetDatabase(c)
	us := &services.User{DB: db}

	keys := &KeyRequest{}
	err := c.BindJSON(&keys)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	err = us.DeletePublicKey(user.ID, keys.Key)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

func getUserFromContext(c *gin.Context) *models.User {
	return c.MustGet("user").(*models.User)
}

type ValidateOtpRequest struct {
	Token string `json:"token"`
}

type EnrollPasskeyRequest struct {
	Name string `json:"name"`
}

type KeyRequest struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
