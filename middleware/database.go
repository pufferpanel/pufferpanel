package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/database"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"gorm.io/gorm"
	"net/http"
)

func NeedsDatabase(c *gin.Context) {
	db, err := database.GetConnection()

	if err != nil {
		logging.Error.Printf("Database not available: %s", err)
		err = pufferpanel.ErrDatabaseNotAvailable
	}

	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Set("db", db)
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

func HasTransaction(c *gin.Context) {
	db := GetDatabase(c)

	if db == nil {
		NeedsDatabase(c)
		db = GetDatabase(c)
		if db == nil {
			response.HandleError(c, pufferpanel.ErrDatabaseNotAvailable, http.StatusInternalServerError)
			return
		}
	}

	c.Set("noTransactionDb", db)

	_ = db.Transaction(func(trans *gorm.DB) error {
		c.Set("db", trans)

		c.Next()

		if c.Errors != nil {
			logging.Error.Printf("Error: %s", c.Errors)
			return errors.New("error in transaction")
		} else if c.Writer.Status() >= 400 {
			return fmt.Errorf("bad status code %d", c.Writer.Status())
		}
		return nil
	})
}
