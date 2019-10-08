package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/apufferi/v3/response"
	"github.com/pufferpanel/pufferpanel/v2"
)

func HasTransaction(c *gin.Context) {
	db := GetDatabase(c)

	if db == nil {
		NeedsDatabase(c)
		db = GetDatabase(c)
		if db == nil {
			response.HandleError(response.From(c), pufferpanel.ErrDatabaseNotAvailable)
			return
		}
	}

	db = db.Begin()

	c.Set("db", db)

	c.Next()

	GetDatabase(c).RollbackUnlessCommitted()
}
