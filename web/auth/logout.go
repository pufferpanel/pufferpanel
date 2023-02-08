package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3/middleware/panelmiddleware"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"net/http"
)

func LogoutPost(c *gin.Context) {
	db := panelmiddleware.GetDatabase(c)
	ss := services.Session{DB: db}

	cookie, err := c.Cookie("puffer_auth")
	if err == http.ErrNoCookie {
		cookie = c.Query("token")
		if cookie != "" {
			err = nil
		}
	}
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	_ = ss.Expire(cookie)

	secure := false
	if c.Request.TLS != nil {
		secure = true
	}

	c.SetCookie("puffer_auth", "", 0, "/", "", secure, true)
	c.SetCookie("puffer_auth_expires", "", 0, "/", "", secure, false)
	c.Status(http.StatusNoContent)
}
