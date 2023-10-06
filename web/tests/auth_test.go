package tests

import (
	"encoding/json"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/database"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"github.com/pufferpanel/pufferpanel/v3/web/auth"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	t.Parallel()
	t.Run("GoodLoginButNoScope", func(t *testing.T) {
		t.Parallel()
		response := CallAPI("POST", "/auth/login", auth.LoginRequestData{
			Email:    loginNoLoginUser.Email,
			Password: loginNoLoginUserPassword,
		}, "")
		assert.Equal(t, http.StatusForbidden, response.Code)
		//ensure we sent back correct headers
		assert.Empty(t, response.Header().Values("Set-Cookie"))
	})
	t.Run("GoodLoginWithLoginScope", func(t *testing.T) {
		t.Parallel()
		response := CallAPI("POST", "/auth/login", auth.LoginRequestData{
			Email:    loginNoServerViewUser.Email,
			Password: loginNoServerViewUserPassword,
		}, "")
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}
		res := &auth.LoginResponse{}
		err := json.NewDecoder(response.Body).Decode(res)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.Equal(t, []*pufferpanel.Scope{pufferpanel.ScopeLogin}, res.Scopes) {
			return
		}
		//ensure we sent back correct headers
		cookies := response.Header().Values("Set-Cookie")
		if !assert.NotEmpty(t, cookies) {
			return
		}
		valid := false
		for _, v := range cookies {
			if strings.HasPrefix(v, "puffer_auth") {
				valid = true
			}
		}
		assert.True(t, valid, "No puffer_auth cookie found")
	})
	t.Run("GoodLoginWithAdminScope", func(t *testing.T) {
		t.Parallel()
		response := CallAPI("POST", "/auth/login", auth.LoginRequestData{
			Email:    loginAdminUser.Email,
			Password: loginAdminUserPassword,
		}, "")
		if !assert.Equal(t, http.StatusOK, response.Code) {
			return
		}
		res := &auth.LoginResponse{}
		err := json.NewDecoder(response.Body).Decode(res)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.Equal(t, []*pufferpanel.Scope{pufferpanel.ScopeAdmin}, res.Scopes) {
			return
		}
		//ensure we sent back correct headers
		tmpRes := http.Response{Header: response.Header()}
		cookies := tmpRes.Cookies()
		if !assert.NotEmpty(t, cookies) {
			return
		}
		valid := false
		for _, v := range cookies {
			if v.Name == "puffer_auth" {
				valid = true
			}
		}
		assert.True(t, valid, "No puffer_auth cookie found")
	})
	t.Run("NoDataLogin", func(t *testing.T) {
		t.Parallel()
		response := CallAPI("POST", "/auth/login", auth.LoginRequestData{
			Email:    "",
			Password: "",
		}, "")
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Empty(t, response.Header().Values("Set-Cookie"))
	})
	t.Run("InvalidEmail", func(t *testing.T) {
		t.Parallel()
		response := CallAPI("POST", "/auth/login", auth.LoginRequestData{
			Email:    "test@notreal.com",
			Password: "testing123",
		}, "")
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Empty(t, response.Header().Values("Set-Cookie"))
	})
	t.Run("InvalidPassword", func(t *testing.T) {
		t.Parallel()
		response := CallAPI("POST", "/auth/login", auth.LoginRequestData{
			Email:    "test@example.com",
			Password: "testing",
		}, "")
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Empty(t, response.Header().Values("Set-Cookie"))
	})
}

func TestLogout(t *testing.T) {
	t.Parallel()
	t.Run("ValidSessionCookie", func(t *testing.T) {
		t.Parallel()
		db, err := database.GetConnection()
		if !assert.NoError(t, err) {
			return
		}
		session, err := createSession(db, loginAdminUser)
		if !assert.NoError(t, err) {
			return
		}
		hashed, err := services.HashToken(session)
		if !assert.NoError(t, err) {
			return
		}

		request, _ := http.NewRequest("POST", "/auth/logout", nil)
		request.AddCookie(&http.Cookie{
			Name:  "puffer_auth",
			Value: session,
		})
		writer := httptest.NewRecorder()
		pufferpanel.Engine.ServeHTTP(writer, request)
		assert.Equal(t, http.StatusNoContent, writer.Code)

		//check to make sure session is gone
		mo := &models.Session{
			Token: hashed,
		}
		var count int64
		err = db.Model(mo).Where(mo).Count(&count).Error
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, int64(0), count)
	})

	t.Run("ValidSessionToken", func(t *testing.T) {
		t.Parallel()
		db, err := database.GetConnection()
		if !assert.NoError(t, err) {
			return
		}
		session, err := createSession(db, loginNoAdminWithServersUser)
		if !assert.NoError(t, err) {
			return
		}
		hashed, err := services.HashToken(session)
		if !assert.NoError(t, err) {
			return
		}

		adminSession, err := createSession(db, loginAdminUser)
		if !assert.NoError(t, err) {
			return
		}
		adminHashed, err := services.HashToken(adminSession)
		if !assert.NoError(t, err) {
			return
		}

		request, _ := http.NewRequest("POST", "/auth/logout?token="+session, nil)
		request.AddCookie(&http.Cookie{
			Name:  "puffer_auth",
			Value: adminSession,
		})
		writer := httptest.NewRecorder()
		pufferpanel.Engine.ServeHTTP(writer, request)
		assert.Equal(t, http.StatusNoContent, writer.Code)

		//check to make sure session is gone
		mo := &models.Session{
			Token: hashed,
		}
		var count int64
		err = db.Model(mo).Where(mo).Count(&count).Error
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, int64(0), count, "Token session not expired")

		mo = &models.Session{
			Token: adminHashed,
		}
		err = db.Model(mo).Where(mo).Count(&count).Error
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, int64(1), count, "Token session expired wrong session")
	})

	t.Run("InvalidSessions", func(t *testing.T) {
		t.Parallel()
		db, err := database.GetConnection()
		if !assert.NoError(t, err) {
			return
		}
		session, err := createSession(db, loginAdminUser)
		if !assert.NoError(t, err) {
			return
		}
		hashed, err := services.HashToken(session)
		if !assert.NoError(t, err) {
			return
		}

		request, _ := http.NewRequest("POST", "/auth/logout", nil)
		request.AddCookie(&http.Cookie{
			Name:  "puffer_auth",
			Value: session + "-extratokens",
		})
		writer := httptest.NewRecorder()
		pufferpanel.Engine.ServeHTTP(writer, request)
		assert.Equal(t, http.StatusNoContent, writer.Code)

		mo := &models.Session{
			Token: hashed,
		}
		var count int64
		err = db.Model(mo).Where(mo).Count(&count).Error
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, int64(1), count)
	})
}

func TestReauth(t *testing.T) {
	t.Parallel()
	t.Run("ReauthNoSession", func(t *testing.T) {
		t.Parallel()
		request, _ := http.NewRequest("POST", "/auth/reauth", nil)
		writer := httptest.NewRecorder()
		pufferpanel.Engine.ServeHTTP(writer, request)
		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	})
	t.Run("ReauthWithValidSession", func(t *testing.T) {
		t.Parallel()
		db, err := database.GetConnection()
		if !assert.NoError(t, err) {
			return
		}
		session, err := createSession(db, loginAdminUser)
		if !assert.NoError(t, err) {
			return
		}

		request, _ := http.NewRequest("POST", "/auth/reauth", nil)
		request.AddCookie(&http.Cookie{
			Name:  "puffer_auth",
			Value: session,
		})
		writer := httptest.NewRecorder()
		pufferpanel.Engine.ServeHTTP(writer, request)
		if !assert.Equal(t, http.StatusOK, writer.Code) {
			return
		}

		response := &auth.LoginResponse{}
		err = json.NewDecoder(writer.Body).Decode(response)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NotEmpty(t, response.Scopes) {
			return
		}

		//ensure we sent back correct headers
		var cookie string
		res := http.Response{Header: writer.Header()}
		cookies := res.Cookies()
		if !assert.NotEmpty(t, cookies) {
			return
		}
		for _, v := range cookies {
			if v.Name == "puffer_auth" {
				cookie = v.Value
			}
		}
		assert.NotEmpty(t, cookie, "No puffer_auth cookie found")

		hashed, err := services.HashToken(cookie)
		if !assert.NoError(t, err) {
			return
		}

		mo := &models.Session{
			Token: hashed,
		}
		var count int64
		err = db.Model(mo).Where(mo).Count(&count).Error
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, int64(1), count)

	})
	t.Run("ReauthWithExpiredSession", func(t *testing.T) {
		t.Parallel()
		db, err := database.GetConnection()
		if !assert.NoError(t, err) {
			return
		}
		session, err := createSession(db, loginAdminUser)
		if !assert.NoError(t, err) {
			return
		}
		hashed, err := services.HashToken(session)
		if !assert.NoError(t, err) {
			return
		}
		err = db.Model(&models.Session{}).Where(&models.Session{Token: hashed}).Updates(&models.Session{ExpirationTime: time.Now().Add(time.Hour * -24)}).Error
		if !assert.NoError(t, err) {
			return
		}

		request, _ := http.NewRequest("POST", "/auth/reauth", nil)
		request.AddCookie(&http.Cookie{
			Name:  "puffer_auth",
			Value: session,
		})
		writer := httptest.NewRecorder()
		pufferpanel.Engine.ServeHTTP(writer, request)
		assert.Equal(t, http.StatusForbidden, writer.Code)
	})
	t.Run("ReauthWithInvalidSession", func(t *testing.T) {
		t.Parallel()
		request, _ := http.NewRequest("POST", "/auth/reauth", nil)
		request.AddCookie(&http.Cookie{
			Name:  "puffer_auth",
			Value: "invalidsession",
		})
		writer := httptest.NewRecorder()
		pufferpanel.Engine.ServeHTTP(writer, request)
		assert.Equal(t, http.StatusForbidden, writer.Code)
	})
}
