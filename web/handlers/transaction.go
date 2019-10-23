package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/v4/response"
	"github.com/pufferpanel/pufferpanel/v2"
	"net/http"
)

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

	db = db.Begin()

	c.Set("db", db)

	c.Next()

	GetDatabase(c).RollbackUnlessCommitted()
}
