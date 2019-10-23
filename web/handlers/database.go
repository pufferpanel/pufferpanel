package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/apufferi/v4/logging"
	"github.com/pufferpanel/apufferi/v4/response"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"net/http"
)

func NeedsDatabase(c *gin.Context) {
	db, err := database.GetConnection()

	if err != nil {
		logging.Exception("Database not available", err)
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
