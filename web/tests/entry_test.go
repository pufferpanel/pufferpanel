package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/database"
	"github.com/pufferpanel/pufferpanel/v3/web"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	_ = os.Remove("testing.db")
	var exitCode = 1

	config.DatabaseDialect.Set("sqlite3", false)
	config.DatabaseUrl.Set("file:testing.db", false)
	config.DaemonEnabled.Set(true, false)
	config.PanelEnabled.Set(true, false)
	config.DatabaseLoggingEnabled.Set(true, false)

	//open db connection
	db, err := database.GetConnection()
	if err != nil {
		panic(err)
	}
	err = prepareUsers(db)
	if err == nil {
		router := gin.New()
		router.Use(gin.Recovery())
		//router.Use(gin.Logger())
		gin.SetMode(gin.ReleaseMode)
		web.RegisterRoutes(router)
		pufferpanel.Engine = router
		exitCode = m.Run()
		database.Close()
	} else {
		fmt.Printf("Error preparing users: %s", err.Error())
	}

	_ = os.Remove("testing.db")
	_ = os.Remove("cache")
	os.Exit(exitCode)
}

func CallAPI(method, url string, body interface{}, token string) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if token != "" {
		request.Header.Add("Authorization", "Bearer "+token)
	}
	writer := httptest.NewRecorder()
	pufferpanel.Engine.ServeHTTP(writer, request)
	return writer
}
