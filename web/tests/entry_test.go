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
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	_ = os.Remove("testing.db")
	var exitCode = 1

	config.DatabaseDialect.Set("sqlite3", false)
	config.DatabaseUrl.Set("file:testing.db", false)
	config.DaemonEnabled.Set(true, false)
	config.PanelEnabled.Set(true, false)
	config.DatabaseLoggingEnabled.Set(false, false)

	_ = os.Mkdir("servers", 0755)
	_ = os.Mkdir("cache", 0755)
	_ = os.Mkdir("binaries", 0755)

	newPath := os.Getenv("PATH")
	fullPath, _ := filepath.Abs(config.BinariesFolder.Value())
	if !strings.Contains(newPath, fullPath) {
		_ = os.Setenv("PATH", newPath+":"+fullPath)
	}

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
	_ = os.RemoveAll("cache")
	_ = os.RemoveAll("servers")
	_ = os.RemoveAll("binaries")

	os.Exit(exitCode)
}

func CallAPI(method, url string, body interface{}, token string) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	return CallAPIRaw(method, url, requestBody, token)
}

func CallAPIRaw(method, url string, body []byte, token string) *httptest.ResponseRecorder {
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	if token != "" {
		request.Header.Add("Authorization", "Bearer "+token)
	}
	writer := httptest.NewRecorder()
	pufferpanel.Engine.ServeHTTP(writer, request)
	return writer
}
