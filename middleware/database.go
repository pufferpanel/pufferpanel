package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"net/http"
)

func NeedsDatabase(c *gin.Context) {
	db, err := database.GetConnection()

	if err != nil {
		logging.Error().Printf("Database not available: %s", err)
		err = pufferpanel.ErrDatabaseNotAvailable
	}

	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Set("db", db)

	c.Next()
}

func GetDatabase(c *gin.Context) *gorm.DB {
	db, exist := c.Get("db")
	if !exist {
		return nil
	}
	casted, ok := db.(*gorm.DB)
	if !ok {
		return nil
	}
	return casted
}
