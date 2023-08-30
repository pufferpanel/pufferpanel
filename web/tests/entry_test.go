package tests

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/database"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/web"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	config.DatabaseDialect.Set("sqlite3", false)
	config.DatabaseUrl.Set("file:testing.db", false)
	config.DaemonEnabled.Set(true, false)
	config.PanelEnabled.Set(true, false)

	//open db connection
	db, err := database.GetConnection()
	if err != nil {
		panic(err)
	}
	err = prepareUsers(db)

	router := gin.New()
	router.Use(gin.Recovery())
	//router.Use(gin.Logger())
	gin.SetMode(gin.ReleaseMode)
	web.RegisterRoutes(router)
	pufferpanel.Engine = router

	exitCode := m.Run()
	database.Close()
	os.Remove("testing.db")
	os.Exit(exitCode)
}

func prepareUsers(db *gorm.DB) error {
	user := &models.User{
		Username:  "idonthaveperms",
		Email:     "noscope@cage.com",
		OtpActive: false,
	}
	err := user.SetPassword("dontletmein")
	if err != nil {
		return err
	}
	err = db.Create(user).Error
	if err != nil {
		return err
	}

	user = &models.User{
		Username:  "testAPI",
		Email:     "test@example.com",
		OtpActive: false,
	}
	err = user.SetPassword("testing123")
	if err != nil {
		return err
	}
	err = db.Create(user).Error
	if err != nil {
		return err
	}

	perms := &models.Permissions{
		UserId: &user.ID,
		User:   *user,
		Scopes: []*pufferpanel.Scope{pufferpanel.ScopeLogin},
	}
	err = db.Create(perms).Error
	return err
}

func CallAPI(method, url string, body interface{}, isAuthenticatedRequest bool) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if isAuthenticatedRequest {
		request.Header.Add("Authorization", "Bearer "+token())
	}
	writer := httptest.NewRecorder()
	pufferpanel.Engine.ServeHTTP(writer, request)
	return writer
}

func token() string {
	return ""
}
