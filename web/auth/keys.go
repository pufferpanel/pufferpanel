package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
)

func GetUserPublicKeys(c *gin.Context) {
	db := middleware.GetDatabase(c)
	us := services.User{DB: db}

	keys, err := us.GetPublicKeys(c.Param("email"))

	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, keys.ToView())
}
